if [ $# -ne 1 ]; then
    echo "./make.sh src_dir"
    exit
fi
go build -o ./bin/ gopl.io/$1
