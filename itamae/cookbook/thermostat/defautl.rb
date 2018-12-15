token = ENV['circle_token']
build_num = ENV['build_num']
cdir = "/tmp/thermostat.#{build_num}"

base_url = "https://circleci.com/api/v1.1/project"

execute 'mkdir' do
   command "mkdir #{cdir}"
   not_if "test -d #{cdir}"
end

[
"curl #{base_url}/github/yuu/thermostat/#{build_num}/artifacts?circle-token=#{token} > data",
"cat data | grep -o 'https://.*[A-z]' > artifacts.txt",
"<artifacts.txt xargs -P4 -I % wget \"%?circle-token=#{token}\" -O thermostat",
].each do |cmd|
    execute 'download' do
        command cmd
        cwd cdir
    end
end

execute 'systemctl stop thermostat' do
    user 'root'
end

execute 'install' do
    command "cp #{cdir}/thermostat /usr/local/bin/thermostat"
    only_if "test -e #{cdir}/thermostat"
end

execute 'systemctl start thermostat' do
    user 'root'
end
