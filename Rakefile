common = ['common.go']

$exes = {
    'server' => common + ['model.go'], 
    'lunch' => common + ['update.go']
}

def getDep(file)
    exe = File.basename(file, '.6')
    return ["#{exe}.go"] + $exes[exe]
end

rule '.6' => lambda { |file| getDep(file) } do |t|
    sh "6g #{t.prerequisites.join ' '}"
end

$exes.each do |exe, dep|
    rule exe => ['.6'] do |t|
        sh "6l -o #{t.name} #{t.source}"
    end
end

desc "Remove derived files"
task :clean do
    FileList['*.6'].include($exes.keys).existing.each do |file|
        sh "rm #{file}"
    end
end

desc "Build all binaries"
task :build_all => $exes.keys

task :default => [:build_all]

desc "Clean, then build"
task :rebuild => [:clean, :build_all]

desc "Format Go source files"
task :format do
    FileList['*.go'].each do |file|
        sh "gofmt -w #{file}"
    end
end

desc "Clean, format, build, and if successful, commit"
task :commit => [:clean, :format, :build_all] do
    sh "git commit -a"
end
