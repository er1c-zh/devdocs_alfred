#!/bin/sh

mkdir -p ./static/icons/docs/
rm -r ./static/icons/docs/*

for i in devdocs/public/icons/docs/* ;
do
        NAME=`basename $i`
        for j in $i/*.png;
        do
                SUFFIX=`basename $j`
                cp $j ./static/icons/docs/$NAME'_'$SUFFIX
        done
done

