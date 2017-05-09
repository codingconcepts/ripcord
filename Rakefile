require 'rake'

OUTPUT = 'ripcord.exe'

task :run do
    Dir.chdir('cmd\standalone') do
        sh("go build -v -o #{OUTPUT}")
        sh(".\\#{OUTPUT}")
    end
end

task :test do
    sh("go test ./... -v")
end

task :bench do
    sh("go test ./... -bench=Benchmark* -v")
end

task :cover do
    sh("go test --coverprofile=coverage.out")
    sh("go tool cover --html=coverage.out")
end