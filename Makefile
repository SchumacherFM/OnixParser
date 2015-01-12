GOFMT=gofmt -s

GOFILES=$(wildcard *.go **/*.go)

OUTDIR=./xmlFiles/

fmt:
	${GOFMT} -w ${GOFILES}

run1:
	rm -f ${OUTDIR}/gonix_*
	go run main.go -infile xmlFiles/pearson-physical-005.xml -db test2 -v -moc 30 -outdir ${OUTDIR}

run2:
	rm -f ${OUTDIR}/gonix_*
	go run main.go -infile xmlFiles/pearson-physical-005.xml -db test2 -v -moc 30 -outdir ${OUTDIR} -logfile zimport.log
