GOFILES=$(wildcard *.go **/*.go)

OUTDIR=./xmlFiles/

fmt:
	gofmt -s -w ${GOFILES}

imp:
	goimports -w ${GOFILES}

run1:
	rm -f ${OUTDIR}/gonix_*
	go run main.go -exitCode -infile xmlFiles/pearson-physical-005.xml -db test \
	-logfile xmlFiles/run.log -tablePrefix gonix_pearson_ -v -outdir ${OUTDIR}

run2:
	rm -f ${OUTDIR}/gonix_*
	go run main.go -infile xmlFiles/pearson-physical-005.xml -db test2 -v -outdir ${OUTDIR} -logfile zimport.log
run:
	go run main.go \
	-db=[dbname] \
	-host=localhost \
	-user=[dbuser] \
	-pass=[dbpassword] \
	-infile=au_wiley_full_20181127.xml \
	-outdir=./xmlFiles/ \
	-tablePrefix=erp_purchase_import_johnwiley_temp_ \
	-logfile=./zipimport.log