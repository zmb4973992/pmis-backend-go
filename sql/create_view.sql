-- 在系统重启后会自动加载，没有使用viper，因为viper不支持sql；使用的是ioutil.readfile

if object_id('disassembly_view') is null
    begin
        execute ('create view disassembly_view as select * from disassembly')
    end

