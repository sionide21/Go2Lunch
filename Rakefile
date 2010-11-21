common = [
  'common.go',
  'model.go',
  'vectors.go',
  'personvector.go',
  'placevector.go'
]

$exes = {
    'lunchd' => common + ['lunchTracker.go'], 
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
    FileList['*.6'].include($exes.keys).include(%w(person place).map { |p| "#{p}vector.go" }).existing.each do |file|
        sh "rm #{file}"
    end
end

desc "Build all binaries"
task :build_all => $exes.keys

task :default => [:generate, :build_all]

desc "Clean, then build"
task :rebuild => [:clean, :build_all]

desc "Format Go source files"
task :format do
    FileList['*.go'].each do |file|
        sh "gofmt -w #{file}"
    end
end

desc "Clean, format, build, and if successful, commit"
task :commit => [:clean, :generate, :format, :build_all] do
    sh "git commit -a"
end

desc "Generates the proper vector files"
task :generate do
  %w(Person Place).each do |type|
    sh "cat #{ENV['GOROOT'] || '~/go'}/src/pkg/container/vector/vector.go | gofmt -r='Vector -> #{type}Vector' | gofmt -r='interface{} -> *#{type}' | sed 's/vector/main/' > #{type.downcase}vector.go"
  end
end
