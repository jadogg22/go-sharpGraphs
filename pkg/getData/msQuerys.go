package getdata

import (
	"fmt"
	"time"
)

// This query is used to collect the report data for the transportation dispachers
// Dispacher | Total Stops | Total Service Failures | Orders with Service Failures | Total Orders // | Total Bill Distance | Total Move Distance | Total Unique Trucks
func MakeTransportationDailyOpsQuery(startStr, endStr string) string {
	return fmt.Sprintf(`
SELECT
    move1.dispatcher_user_id,
    SUM(COALESCE(stop_counts.stop_count, 0)) AS total_stops,
    SUM(COALESCE(servicefail_counts.servicefail_count, 0)) AS total_servicefail_count,
    COUNT(DISTINCT CASE WHEN servicefail_counts.servicefail_count > 0 THEN mo1.order_id END) AS orders_with_service_fail,
    COUNT(DISTINCT mo1.order_id) AS total_orders,
    SUM(COALESCE(p.empty_distance, 0)) AS total_empty_distance,  -- Replaced bill_distance with empty_distance
    SUM(COALESCE(p.loaded_distance, 0)) AS total_loaded_distance,  -- Replaced move_distance with loaded_distance
    COUNT(DISTINCT continuity.equipment_id) AS total_unique_trucks
FROM
    movement move1
    JOIN movement_order mo1 ON move1.id = mo1.movement_id AND mo1.company_id = 'TMS'
    JOIN orders ON orders.id = mo1.order_id AND orders.company_id = 'TMS'
    JOIN prorated_orderdist p ON orders.id = p.order_id  -- Join prorated_orderdist to get empty and loaded distances
    LEFT JOIN (
        SELECT
            orders.id AS order_id,
            COUNT(stop.id) AS stop_count
        FROM
            orders
            JOIN stop ON stop.order_id = orders.id AND stop.company_id = 'TMS' AND stop.status <> 'V'
        WHERE
            orders.company_id = 'TMS' AND orders.status <> 'V'
        GROUP BY
            orders.id
    ) AS stop_counts ON stop_counts.order_id = mo1.order_id
    LEFT JOIN (
        SELECT
            servicefail.order_id,
            COUNT(DISTINCT servicefail.id) AS servicefail_count
        FROM
            servicefail
        WHERE
            servicefail.company_id = 'TMS'
            AND servicefail.status <> 'V'
        GROUP BY
            servicefail.order_id
    ) AS servicefail_counts ON servicefail_counts.order_id = mo1.order_id
    LEFT JOIN continuity ON continuity.movement_id = move1.id
        AND continuity.equipment_type_id = 'T'
WHERE
    move1.company_id = 'TMS'
    AND move1.loaded = 'L'
    AND move1.status <> 'V'
    AND move1.id IN (
        SELECT
            movement.id
        FROM
            movement
            JOIN movement_order ON movement.id = movement_order.movement_id AND movement_order.company_id = 'TMS'
            JOIN orders ON orders.id = movement_order.order_id AND orders.company_id = 'TMS'
            JOIN stop ON stop.id = orders.shipper_stop_id AND stop.company_id = 'TMS'
            JOIN stop dest ON dest.id = orders.consignee_stop_id AND dest.company_id = 'TMS'
        WHERE
            movement.company_id = 'TMS'
            AND stop.actual_arrival BETWEEN {ts '%s 00:00:00'} AND {ts '%s 23:59:59'}
            AND movement.status <> 'V'
            AND movement.loaded = 'L'
    )
GROUP BY
    move1.dispatcher_user_id
ORDER BY
    move1.dispatcher_user_id;
`, startStr, endStr)
}

func getLogisticsMTDQuery(startDate, endDate time.Time) string {
	startdateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	query := fmt.Sprintf(`
SELECT
    move1.dispatcher_user_id,
    SUM(orders.total_charge) AS revenue,           -- Total charges per dispatcher
    SUM(move1.override_pay_amt) AS override_pay_amt,  -- Sum of override pay amounts
    SUM(COALESCE(de_sum.amount, 0) + override_pay_amt) AS truck_hire,  -- Total truck hire (sum of driver_extra_pay amounts + override pay)
    SUM(COALESCE(stop_counts.stop_count, 0)) AS total_stops,  -- Total stops per dispatcher
    SUM(COALESCE(servicefail_counts.servicefail_count, 0)) AS total_servicefail_count,  -- Total service failures per dispatcher
    COUNT(DISTINCT CASE WHEN servicefail_counts.servicefail_count > 0 THEN mo1.order_id END) AS orders_with_service_fail,  -- Count of orders with service fails
    COUNT(DISTINCT mo1.order_id) AS total_orders,  -- Count of distinct orders per dispatcher
	SUM(orders.bill_distance) AS total_bill_distance  -- Sum of bill distances per dispatcher


FROM
    movement move1
    JOIN movement_order mo1 ON move1.id = mo1.movement_id AND mo1.company_id = 'TMS2'
    JOIN orders ON orders.id = mo1.order_id AND orders.company_id = 'TMS2'
    LEFT JOIN (
        SELECT
            de.movement_id,
            SUM(de.amount) AS amount
        FROM
            driver_extra_pay de
        WHERE
            de.company_id = 'TMS2'
        GROUP BY
            de.movement_id
    ) AS de_sum ON de_sum.movement_id = move1.id
    LEFT JOIN (
        SELECT
            orders.id AS order_id,
            COUNT(stop.id) AS stop_count
        FROM
            orders
            JOIN stop ON stop.order_id = orders.id AND stop.company_id = 'TMS2' AND stop.status <> 'V'
        WHERE
            orders.company_id = 'TMS2' AND orders.status <> 'V'
        GROUP BY
            orders.id
    ) AS stop_counts ON stop_counts.order_id = mo1.order_id
    LEFT JOIN (
        SELECT
            servicefail.order_id,
            COUNT(DISTINCT servicefail.id) AS servicefail_count
        FROM
            servicefail
        WHERE
            servicefail.company_id = 'TMS2'
            AND servicefail.status <> 'V'
        GROUP BY
            servicefail.order_id
    ) AS servicefail_counts ON servicefail_counts.order_id = mo1.order_id
WHERE
    move1.company_id = 'TMS2'
    AND move1.loaded = 'L'
    AND move1.status <> 'V'
    AND move1.id IN (
        SELECT
            movement.id
        FROM
            movement
            LEFT OUTER JOIN payee ON payee.id = movement.override_payee_id AND payee.company_id = 'TMS2'
            LEFT OUTER JOIN users ON users.id = movement.dispatcher_user_id AND users.company_id = 'TMS2'
            JOIN movement_order ON movement.id = movement_order.movement_id AND movement_order.company_id = 'TMS2'
            JOIN orders ON orders.id = movement_order.order_id AND orders.company_id = 'TMS2'
            LEFT OUTER JOIN customer ON customer.id = orders.customer_id AND customer.company_id = 'TMS2'
            LEFT OUTER JOIN revenue_code ON revenue_code.id = orders.revenue_code_id AND revenue_code.company_id = 'TMS2'
            JOIN stop ON stop.id = orders.shipper_stop_id AND stop.company_id = 'TMS2'
            JOIN stop dest ON dest.id = orders.consignee_stop_id AND dest.company_id = 'TMS2'
        WHERE
            movement.company_id = 'TMS2'
            AND stop.actual_arrival BETWEEN {ts '%s 00:00:00'} AND {ts '%s 23:59:59'}
            AND movement.status <> 'V'
            AND movement.loaded = 'L'
    )
GROUP BY
    move1.dispatcher_user_id
ORDER BY
    move1.dispatcher_user_id;
`, startdateStr, endDateStr)

	return query
}

func MakeQuery(yesterday string) string {

	return fmt.Sprintf(`SELECT 
    movement.id AS move_id,
    movement.move_distance AS move_distance,
    movement.loaded AS loaded,
    orders.id AS order_id,
    orders.total_charge AS charges,
    orders.bill_distance AS bill_distance,
    orders.freight_charge AS freight_charge,
    origin.city_name AS origin_city,
    origin.state AS origin_state,
    continuity.equipment_id AS equip_id,
    continuity.actual_arrival AS actual_arrival,
    continuity.dest_actualarrival AS del_date,
    continuity.equipment_id AS tractor,
    continuity.equipment_type_id,
    tractor.dispatcher AS dispatcher,
    tractor.fleet_id AS fleet_id,
    fleet.description AS fleet_description,
    users.name AS user_name,
    COUNT(servicefail.id) AS servicefail_count,
    CASE WHEN COUNT(servicefail.id) > 0 THEN 1 ELSE 0 END AS has_servicefail,
	stop_count.stop_count AS stop_count
FROM 
    movement
    JOIN movement_order ON movement.id = movement_order.movement_id AND movement_order.company_id = 'TMS'
    JOIN orders ON orders.id = movement_order.order_id AND orders.company_id = 'TMS'
    JOIN stop origin ON origin.movement_id = movement.id AND origin.movement_sequence = 1 AND origin.company_id = 'TMS'
    JOIN stop dest ON dest.id = movement.dest_stop_id AND dest.company_id = 'TMS'
 JOIN continuity ON movement.id = continuity.movement_id AND continuity.equipment_type_id = 'T' AND continuity.company_id = 'TMS'
    JOIN tractor ON tractor.id = continuity.equipment_id AND tractor.company_id = 'TMS'
    LEFT JOIN fleet ON fleet.id = tractor.fleet_id AND fleet.company_id = 'TMS'
    LEFT JOIN users ON users.id = tractor.dispatcher AND users.company_id = 'TMS'
    LEFT JOIN servicefail ON servicefail.order_id = orders.id AND servicefail.company_id = 'TMS'
	LEFT JOIN (
        SELECT order_id, COUNT(*) AS stop_count
        FROM stop
        WHERE company_id = 'TMS'
        GROUP BY order_id
    ) AS stop_count ON stop_count.order_id = orders.id

WHERE 
    movement.company_id = 'TMS' 
    AND continuity.dest_actualarrival >= {ts '%s 00:00:00'}
    AND continuity.dest_actualarrival <= {ts '%s 23:59:59'}
    AND movement.status <> 'V'
GROUP BY 
    movement.id, 
    movement.move_distance, 
    movement.loaded, 
    orders.id,
    orders.total_charge, 
    orders.bill_distance, 
    orders.freight_charge, 
    origin.city_name, 
    origin.state,
    continuity.equipment_id,
    continuity.actual_arrival, 
    continuity.dest_actualarrival, 
    continuity.equipment_id, 
    continuity.equipment_type_id,
    tractor.dispatcher, 
    tractor.fleet_id, 
    fleet.description, 
    users.name, 
	stop_count.stop_count
ORDER BY 
    dispatcher, 
    tractor, 
    continuity.dest_actualarrival;`, yesterday, yesterday)
}

// The query is old and not used but I w ant to keep it for reference becaue the query is complex and
// has some good bits I might need to come and reference later

func OldTrandportaionOrdersQuery(startDateStr, endDateStr string) string {
	return fmt.Sprintf(`select 
			orders.id order_id, orders.operations_user operations_user, orders.revenue_code_id revenue_code_id,
			orders.freight_charge freight_charge, orders.bill_distance bill_miles, orders.bill_date bill_date,
			orders.ctrl_party_id controlling_party, orders.commodity_id commodity, orders.order_type_id order_type,
			orders.equipment_type_id order_trailer_type, origin.state origin_value, dest.state destination_value, customer.id customer_id,
			customer.name customer_name, customer.category customer_category, category.descr category_descr, movement.id movement_id,
			loaded, move_distance, movement.brokerage, trailer.trailer_type trailer_type, origin.city_name origin_city, origin.state origin_state,
			dest.city_name dest_city, dest.state dest_state, other_charge.amount oc_amount, charge_code.is_fuel_surcharge is_fuel_surcharge,
			dest.sched_arrive_early report_date, dest.actual_arrival actual_date, prorated_orderdist.empty_distance empty_miles,
			prorated_orderdist.loaded_distance loaded_miles, (prorated_orderdist.empty_distance+prorated_orderdist.loaded_distance) total_miles,
			orders.id record_count, orders.id fuel_surcharge, orders.id remaining_charges, orders.id total_revenue, orders.id empty_pct,
			orders.id rev_loaded_mile, orders.id rev_total_mile, orders.id billed, orders.id week_value, orders.id month_value,
			orders.id quarter_value, revenue_code_id detail_id 
		from 
			orders left outer join customer on customer.id = orders.customer_id and customer.company_id = 'TMS'  
			left outer join category on category.id = customer.category and category.company_id = 'TMS'  
			left outer join movement_order on movement_order.order_id = orders.id and movement_order.company_id = 'TMS'  
			left outer join movement on movement.id = movement_order.movement_id and movement.company_id = 'TMS'  
			left outer join continuity trailercont on (movement.id = trailercont.movement_id)and(trailercont.equipment_type_id='L') and  trailercont.company_id = 'TMS'  
			left outer join trailer on trailer.id = trailercont.equipment_id and trailer.company_id = 'TMS'  left outer join other_charge on other_charge.order_id = orders.id and other_charge.company_id = 'TMS'  
			left outer join charge_code on charge_code.id = other_charge.charge_id  
			left outer join prorated_orderdist on prorated_orderdist.order_id = orders.id and prorated_orderdist.company_id = 'TMS'  ,stop origin ,stop dest 
		where 
			orders.company_id = 'TMS' and orders.status <> 'Q' and orders.status <> 'V' and (orders.subject_order_status is null or orders.subject_order_status <> 'S') and loaded = 'L' 
			and ((dest.actual_arrival is not null and dest.actual_arrival >= {ts '2024-08-11 00:00:00'}) or dest.actual_arrival is null and dest.sched_arrive_early >= {ts '%s 00:00:00'}) and ((dest.actual_arrival is not null and dest.actual_arrival <= {ts '%s 23:59:59'}) 
			or dest.actual_arrival is null and dest.sched_arrive_early <= {ts '2024-08-15 23:59:59'}) and origin.id = orders.shipper_stop_id  and  origin.company_id = 'TMS' and dest.id = orders.consignee_stop_id  and  dest.company_id = 'TMS' order by revenue_code_id, order_id, movement_id`, startDateStr, endDateStr)
}

func GetVacationHoursByCompanyQuery(companyID string) string {
	if companyID == "drivers" {
		return DriversVacationHoursQuery()
	}

	return fmt.Sprintf(`WITH LatestLeaveTransaction AS (
    SELECT
        payee_id,
        MAX(trx_date) AS latest_trx_date
    FROM leave_transaction
    WHERE company_id = '%s'
      AND applies_to = 'V'
      AND (effect = 'B' OR effect = 'S')
      AND (is_void IS NULL OR is_void <> 'Y')
    GROUP BY payee_id
),
NewestLeaveTransaction AS (
    SELECT
        lt.payee_id,
        lt.amount,
        lt.trx_date
    FROM leave_transaction lt
    JOIN LatestLeaveTransaction llt
      ON lt.payee_id = llt.payee_id
     AND lt.trx_date = llt.latest_trx_date
    WHERE lt.company_id = '%s'
      AND lt.applies_to = 'V'
      AND (lt.effect = 'B' OR lt.effect = 'S')
      AND (lt.is_void IS NULL OR lt.is_void <> 'Y')
),
OfficeEmployees AS (
    SELECT
        payee.*,
        off_payee.regular_rate,
        off_payee.vacation_hours_due,
        off_payee.vacation_pay_rate
    FROM
        payee
        JOIN off_payee 
          ON payee.id = off_payee.id 
         AND payee.company_id = off_payee.company_id
        JOIN drs_payee 
          ON payee.id = drs_payee.id 
         AND payee.company_id = drs_payee.company_id
    WHERE
        payee.company_id = '%s'
        AND payee.office_employee = 'Y'
        AND off_payee.company_id = '%s'
        AND status = 'A'
)
SELECT
    --oe.*,
	oe.id,
	oe.check_name,
	oe.vacation_pay_rate,
    nlt.amount AS latest_amount
FROM
    OfficeEmployees oe
    LEFT JOIN NewestLeaveTransaction nlt
      ON oe.id = nlt.payee_id
ORDER BY
    oe.id;`, companyID, companyID, companyID, companyID)
}

func DriversVacationHoursQuery() string {
	return `
WITH LatestTransactions AS (
    SELECT 
        payee_id,
        amount,
        trx_date,
        ROW_NUMBER() OVER (PARTITION BY payee_id ORDER BY trx_date DESC, id DESC) AS rn
    FROM 
        leave_transaction
    WHERE 
        company_id = 'TMS' 
        AND applies_to = 'V' 
        AND (is_void IS NULL OR is_void <> 'Y') 
        AND effect <> 'S'
)

SELECT 
    p.id, 
    p.check_name, 
    d.vacation_pay_rate, 
    lt.amount as latest_amout
FROM 
    payee p
JOIN 
    drs_payee d ON p.id = d.id 
LEFT JOIN 
    LatestTransactions lt ON p.id = lt.payee_id AND lt.rn = 1
WHERE 
    p.company_id = 'TMS' 
    AND p.non_office_emp = 'Y' 
    AND p.status = 'A'
    AND d.company_id = 'TMS';
`
}

func MakeCodedRevenueQuery(startDate, endDate time.Time) string {
	startdateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	return fmt.Sprintf(`SELECT revenue_code_id, freight_charge FROM orders where bol_recv_date between {ts '%s 00:00:00'} and {ts '%s 23:59:59'} and company_id = 'TMS'`, startdateStr, endDateStr)
}

func MakeStackedMilesQuery(startDate, endDate time.Time) string {
	startdateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	return fmt.Sprintf(`SELECT o.id, o.bol_recv_date, p.loaded_distance, p.empty_distance from orders o join prorated_orderdist p on o.id = p.order_id where o.bol_recv_date between {ts '%s 00:00:00'} and {ts '%s 23:59:59'} and o.company_id = 'TMS'`, startdateStr, endDateStr)
}

// New query for the sportsmans loads report with price breakdown and pallets
// Mcloud has a database bug that causes the pallets_picked_up and pallets_dropped to be switched

// If they ever fix this bug, swap the two fields in the query and remove the AS statements
func MakeSportsmansQuery(startDate, endDate string) string {

	return fmt.Sprintf(`SELECT 
    o.id AS order_id,
    o.ordered_date,
    s.actual_arrival AS DEL_DATE,
    o.bill_date,
    s.city_name,
    s.state,
    s.zip_code,
    s.location_name AS consignee,
    o.bill_distance AS miles,
    o.blnum AS bol_number,
    o.commodity,
    s.weight,
    s.movement_sequence,

    -- mcloud database bug
    s.pallets_picked_up AS pallets_dropped,
    s.pallets_dropped AS pallets_picked_up,

    o.freight_charge,
    o.otherchargetotal,
    o.total_charge,
    
    -- Separate charge sums for FUD, EDR, EPU
    SUM(CASE WHEN oc.charge_id = 'FUD' THEN COALESCE(oc.amount, 0) ELSE 0 END) AS fuel_surcharge,
    SUM(CASE WHEN oc.charge_id = 'EDR' THEN COALESCE(oc.amount, 0) ELSE 0 END) AS extra_drops,
    SUM(CASE WHEN oc.charge_id = 'EPU' THEN COALESCE(oc.amount, 0) ELSE 0 END) AS extra_pickup,

    -- Calculate 'other_charge' as the residual of otherchargetotal minus the sum of FUD, EDR, and EPU
    ROUND(o.otherchargetotal 
        - SUM(CASE WHEN oc.charge_id = 'FUD' THEN COALESCE(oc.amount, 0) ELSE 0 END)
        - SUM(CASE WHEN oc.charge_id = 'EDR' THEN COALESCE(oc.amount, 0) ELSE 0 END)
        - SUM(CASE WHEN oc.charge_id = 'EPU' THEN COALESCE(oc.amount, 0) ELSE 0 END), 2) AS other_charge,

    -- Calculate per pallet dropped costs
    ROUND(SUM(CASE WHEN oc.charge_id = 'FUD' THEN COALESCE(oc.amount, 0) ELSE 0 END) / NULLIF(s.pallets_picked_up, 0), 2) AS per_pallet_dropped_fuel_surcharge,
    ROUND(o.freight_charge / NULLIF(s.pallets_picked_up, 0), 2) AS per_pallet_dropped_freight_charge,

    -- Combine carrier trailer from movement or equipment_item
    COALESCE(m.carrier_trailer, ei.equipment_id) AS carrier_trailer  -- Use movement first, fallback to equipment_item

FROM 
    orders o
JOIN 
    stop s ON o.id = s.order_id AND o.company_id = s.company_id
LEFT OUTER JOIN 
    other_charge oc ON oc.order_id = o.id AND oc.company_id = 'TMS'
LEFT OUTER JOIN 
    movement m ON o.curr_movement_id = m.id
LEFT OUTER JOIN 
    equipment_item ei ON m.equipment_group_id = ei.equipment_group_id 
    AND ei.company_id = 'TMS' 
    AND ei.equipment_type_id = 'T'

WHERE 
    CAST(o.bill_date AS DATE) BETWEEN '%s' AND '%s'
    AND o.customer_id = 'SPORTSUT'
    AND (oc.charge_id IN ('FUD', 'EDR', 'EPU') OR oc.charge_id IS NULL)
    
GROUP BY 
    o.id, s.actual_arrival, o.ordered_date, o.bill_date, s.city_name, s.state, s.zip_code, 
    s.location_name, o.bill_distance, o.blnum, o.commodity, s.weight, s.movement_sequence, 
    s.pallets_dropped, s.pallets_picked_up, o.freight_charge, o.otherchargetotal, o.total_charge,
    m.carrier_trailer, ei.equipment_id  -- Include both carrier_trailer and equipment_id

ORDER BY 
    o.id, s.movement_sequence, o.bill_date;
`, startDate, endDate)
}

func DashboardQuery(date1, date2 string) string {
	return fmt.Sprintf(`select freight_charge, otherchargetotal, total_charge, xferred2billing, orders.id, city_name, state, actual_arrival, sched_arrive_early, revenue_code_id, bill_distance from orders ,stop 
where orders.company_id = 'TMS' and orders.status <> 'V' and orders.status <> 'Q' and (orders.subject_order_status is null or orders.subject_order_status <> 'S')
and actual_arrival between {ts '%s'} and {ts '%s'}
and stop.id=orders.shipper_stop_id and stop.company_id = 'TMS' order by revenue_code_id, orders.id`, date1, date2)
}

func GetDriverManagerQuery() string {
	return `SELECT 
    COALESCE(d.id, 'Unknown') AS driver_id,         -- Replace NULL driver_id with 'Unknown'
    COALESCE(d.fleet_manager, 'Unassigned') AS fleet_manager, -- Replace NULL fleet_manager with 'Unassigned'
    DATEPART(week, c.actual_arrival) AS week_number,
    DATENAME(MONTH, c.actual_arrival) AS month_name,
    DATEPART(month, c.actual_arrival) AS month_order,
    SUM(m.move_distance) AS total_move_distance
FROM 
    movement AS m
INNER JOIN 
    continuity AS c ON c.movement_id = m.id
INNER JOIN 
    driver AS d ON d.id = c.equipment_id
WHERE 
    d.is_active = 'y'
    AND d.company_id = 'TMS'
    AND c.equipment_type_id = 'd'
    AND c.actual_arrival BETWEEN '2025-01-01' AND GETDATE() + 1
    AND m.status <> 'V'
GROUP BY 
    COALESCE(d.id, 'Unknown'),
    COALESCE(d.fleet_manager, 'Unassigned'),
    DATEPART(week, c.actual_arrival), 
    DATENAME(MONTH, c.actual_arrival), 
    DATEPART(month, c.actual_arrival)
ORDER BY 
    COALESCE(d.id, 'Unknown'),
    month_order, 
    week_number;
`
}
