SELECT DeliveryDate,
	Week as Name,
	strftime('%m', DeliveryDate) as NameStr,
	SUM(LoadedMiles) AS TotalLoadedMiles,
	SUM(EmptyMiles) AS TotalEmptyMiles,
	SUM(TotalMiles) AS TotalMiles,
	SUM(EmptyMiles) / SUM(TotalMiles) * 100 AS PercentEmpty
FROM transportation
WHERE Month BETWEEN "2024 M01" AND "2024 M02"
GROUP BY Name;