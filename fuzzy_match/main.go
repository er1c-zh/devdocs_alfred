package fuzzy_match

import (
	"github.com/er1c-zh/go-now/log"
	"strings"
	"time"
)

type Candidate interface {
	GetString() string
}

type Option struct {
	MinScore int64
	Limit    int
	Debug    bool
}

func NewOption() *Option {
	return &Option{}
}

func (opt *Option) WithMinScore(minScore int64) *Option {
	opt.MinScore = minScore
	return opt
}

func (opt *Option) WithLimit(limit int) *Option {
	opt.Limit = limit
	return opt
}

func (opt *Option) WithDebug() *Option {
	opt.Debug = true
	return opt
}

func FuzzyMatch[T Candidate](pattern string, candidatePool []T, opt *Option) []T {
	t0 := time.Now()
	defer func() {
		log.Info("FuzzyMatch len(candidate) = %d, cost %dms", len(candidatePool), time.Now().Sub(t0).Milliseconds())
	}()
	if opt == nil {
		opt = NewOption()
	}
	if opt.Limit == 0 {
		opt.Limit = len(candidatePool)
	}

	resultList := make([]*MatchResult, opt.Limit)
	for _, candidate := range candidatePool {
		result := fuzzyMatch(pattern, candidate, opt)
		{
			ptr := &result
			for i := 0; i < len(resultList); i += 1 {
				if resultList[i] == nil {
					resultList[i] = ptr
					break
				}
				if ptr.Score > resultList[i].Score {
					ptr, resultList[i] = resultList[i], ptr
				}
			}
		}
	}

	afterMatch := make([]T, 0)
	for _, result := range resultList {
		if result == nil {
			continue
		}
		if result.Score < opt.MinScore {
			continue
		}
		afterMatch = append(afterMatch, result.Candidate.(T))
	}
	return afterMatch
}

type proc struct {
	ip, ic        int64
	score         int64
	consecutively int64
}

func (p proc) Fork() proc {
	return p
}

type MatchResult struct {
	// TODO match index
	Candidate
	Score int64
}

func fuzzyMatch(pattern string, candidate Candidate, opt *Option) MatchResult {
	var r MatchResult
	r.Candidate = candidate
	queue := []proc{
		{
			ip:    0,
			ic:    0,
			score: 0,
		},
	}

	sp, sc := strings.ToLower(pattern), strings.ToLower(candidate.GetString())

	lp, lc := int64(len(pattern)), int64(len(candidate.GetString()))
	for len(queue) != 0 {
		nq := make([]proc, 0)
		for _, p := range queue {
			if p.ip >= lp || p.ic >= lc {
				if r.Score < p.score {
					r.Score = p.score
				}
				continue
			}
			{
				// 匹配
				if sp[p.ip] == sc[p.ic] {
					np := p.Fork()
					np.consecutively += 1
					// 计算分数
					// 0. 基础分数
					np.score += 1
					// 1. 靠前的匹配更多
					if np.ic-np.ip < 5 {
						np.score += 5 - (np.ic - np.ip)
					}
					// 2. 连续匹配分数更高
					np.score += np.consecutively
					// 3. TODO 匹配分隔符后的字符更好
					// 4. TODO 精确匹配大小写的更好

					// 更新下标
					np.ip += 1
					np.ic += 1
					nq = append(nq, np)
				} else {
					// 清除连续匹配的得分
					p.consecutively = 0
				}
			}
			{
				// 不匹配直接跳过
				p.ic += 1
				nq = append(nq, p)
			}
		}
		queue = nq
	}

	if opt.Debug {
		log.Debug("fuzzyMatch pattern: '%s', candidate: '%s', score: %d", sp, sc, r.Score)
	}

	return r
}
