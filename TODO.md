# TODO


- [X] migrate database to mysql on docker - just used a volume for now
- [X] write tests to verify that the database is working
- [x] write tests to verify the backend is working as expected
- [X] get backend working with that mysql docker container
- [X] get backend dockerized
- [X] get frontend dockerized
- [X] create a docker-compose file to run the whole thing
- [X] write script to update the backend from Github - docker down, git pull, docker up

# Done

- [] Finish refactoring the first endpoint YearbyYear
   - [X] refactor for the data structure
   - [X] refactor for the old data
   - [ ] refactor for the new data
- [] Push out the pie chart
- [] Finish Dispatchers
    - [X] Finish the dispatcher parser
    - [] Finish the front end
    - [] Finish backend 

# Notes

the backend is really starting to come together at this point. I got the entier db fully ported over to postgress. Its fully ported to a docker container that will automatically open up when it goes down so I'm pretty stoked about that. I downloaded pgAdmin for some kind of gui and lemme tell you it is not intuitive in the slightest. I'm googling things left and right just trying to find my damn tables. Its all good now i've got a rytham going but I'm really brushing up on my sql skills. If only it was bash or something id be pro. 


## DB migrated 

now that the db is migrated and we have worked out some of the kinks that were getting in our way. we need to move onto the daily ops side of the stats. For now we just have a basic table that doest have any logic corisponding to it. The main goal of this table is to track how well each of the dispachers are doing. Miles, deadhead. etc. Typically luke does this every morning but we are going to automate it. The hardest part of this is getting the data to the computer I think we're going to need somone to click a button until I can get the data from my email.

regardless I think that it should be having some kind of conditional formating and that is simple enough to do. Right now I have a db full of "dispacher" data but its really hard to compair because it needs to almost be on a per truck basis. So we need to change the db to take in the data on a per truck basis. Once we have that information we can start to do some really fun sql queries. The basic one that I wrote to get some of the data that we need is here below. I just need to change it to organize the data by day as well.

Okay today is aug 5th and I've descided that enough is enough. we're going full automation rn. I've got access to the mcloud sql server. and on ssms ive been running quarries to get the data that I need. I've reverse engenereed one of the quaries that it ran and I'm going to use that to get the data that I need

```sql

select 
	movement.id move_id, movement.move_distance move_distance, movement.loaded loaded, orders.id order_id,
	orders.total_charge charges, orders.bill_distance bill_distance, orders.freight_charge freight_charge, origin.city_name origin_city, origin.state origin_state,
	continuity.equipment_id equip_id,
	continuity.actual_arrival actual_arrival, continuity.dest_actualarrival del_date, continuity.equipment_id tractor, continuity.equipment_type_id,
	tractor.dispatcher dispatcher, tractor.fleet_id fleet_id, fleet.description fleet_description, users.name user_name, orders.bill_date bill_date
from 
	movement ,movement_order ,orders ,stop origin ,stop dest ,continuity ,tractor left outer join fleet on fleet.id = tractor.fleet_id and fleet.company_id = 'TMS'  left outer join users on users.id = tractor.dispatcher and users.company_id = 'TMS'  
where 
	movement.company_id = 'TMS' and (continuity.dest_actualarrival >= {ts '2024-08-01 00:00:00'} and continuity.dest_actualarrival <= {ts '2024-08-03 23:59:59'}) and movement.status <> 'V' and movement.id = movement_order.movement_id and movement_order.company_id = 'TMS' and orders.id = movement_order.order_id and orders.company_id = 'TMS' and origin.movement_id = movement.id and origin.movement_sequence = 1 and  origin.company_id = 'TMS' and dest.id = movement.dest_stop_id  and  dest.company_id = 'TMS' and (movement.id = continuity.movement_id)and(continuity.equipment_type_id='T') and continuity.company_id = 'TMS' and tractor.id = continuity.equipment_id and tractor.company_id = 'TMS' 
order by 
	dispatcher, tractor, continuity.dest_actualarrival
```

I've now got the connection working. i'm just going to need to create a struct, for these values. and add them to my database so we don't stress out or cause problems for the server. either way i'm stoked at all the progress that I made Somehow I'm going to need to add the ORDER AND STOP percentages to these records as well.

I'm thinking of doing some kind of join to the table where I get the order ID and for each order ID i pull up the service incidents. With that I can see if there are any service incidents for the order. the Order% can be likly calculated with a booleen and then the stop percentage can be calculated with a count of the stops.

so ive added the following to the sql query
```sql
SELECT 
	count (servicefail.id) as servicefail_count,
	case when count(servicefail.id) > 0 then 1 else 0 end as servicefail
	stop_count.stop_count as stop_count
FROM 
 LEFT JOIN servicefail ON servicefail.order_id = orders.id AND servicefail.company_id = 'TMS'
	LEFT JOIN (
        SELECT order_id, COUNT(*) AS stop_count
        FROM stop
        WHERE company_id = 'TMS'
        GROUP BY order_id
    ) AS stop_count ON stop_count.order_id = orders.id

group by
	stop_count.stop_count
```

This extra data helps me calcuate the stop percentages in the db!




- [X] setup table in my db
- [X] grab old data for compairisons sake
- [X] actually write some frontend code for this data
- [X] figure out the order and stop percentages.

- [X] finish the pie chart



Okay now we are getting somewhere. we really have just ritten and rewritten this very many times. Although its taking me this long we have sor really cool functionality. We have gone from a one off graph to now having a fully automated systems of graphs.

- [] finish Designing the dashboard
- [] finish the backnd for the dashboard

- [] finsih the pie chart - The tool tip sucks. and the spacing sucks.
- [] Maybe include the average order reveneue just to more easlily compare the differances 


