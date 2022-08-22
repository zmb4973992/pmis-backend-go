-- 在系统重启后会自动加载，没有使用viper，因为viper不支持sql；使用的是ioutil.readfile



if object_id('disassembly_view') is null
    begin
        execute ('create view disassembly_view as select * from disassembly')
    end


if object_id('a') is null
    begin
        execute (
'create view a as
with t1 as (select distinct disassembly_id, first_value(actual_progress) over ( partition by disassembly_id order by date desc) as latest_actual_progress
            from work_progress
            where actual_progress is not null
)
select project_id,
       work_progress.disassembly_id  as disassembly_id,
       max(level) as level,
       date,
       max(planned_progress)  as planned_progress,
       max(predicted_progress)    as predicted_progress,
       max(actual_progress)    as actual_progress,
       max(latest_actual_progress) as latest_actual_progress
from work_progress
         left join t1 on work_progress.disassembly_id = t1.disassembly_id
         join disassembly on work_progress.disassembly_id = disassembly.ID
group by project_id, work_progress.disassembly_id, date'
)
    end