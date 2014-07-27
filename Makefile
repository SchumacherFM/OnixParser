GOFMT=gofmt -s

GOFILES=$(wildcard *.go **/*.go)

OUTDIR=./xmlFiles/

format:
	${GOFMT} -w ${GOFILES}

run1:
	rm -f ${OUTDIR}/gonix_*
	go run OnixParser.go -infile xmlFiles/oup_onix.xml -db test2 -v -moc 30 -outdir ${OUTDIR}

run2:
	rm -f ${OUTDIR}/gonix_*
	go run OnixParser.go -infile xmlFiles/oup_onix.xml -db test2 -v -moc 30 -outdir ${OUTDIR} -logfile zimport.log
