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
    SUM(orders.bill_distance) AS total_bill_distance,
    SUM(move1.move_distance) AS total_move_distance,
    COUNT(DISTINCT continuity.equipment_id) AS total_unique_trucks
FROM
    movement move1
    JOIN movement_order mo1 ON move1.id = mo1.movement_id AND mo1.company_id = 'TMS'
    JOIN orders ON orders.id = mo1.order_id AND orders.company_id = 'TMS'
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

func DriversVacationHoursQuery(companyID string) string {
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
        AND payee.non_office_emp = 'Y'
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
