let data = [{"dispacher":"CAMI HANSEN","total_orders":162,"revenue":313371.20000000024,"truck_hire":253652.69999999995,"net_revenue":59718.50000000029,"margins":0.1905679271100862,"total_miles":117625,"rev_per_mile":0.5077024442082916,"stop_percentage":80.2547770700637,"order_percentage":67.90123456790124},{"dispacher":"JERRAMI MAROTZ","total_orders":16,"revenue":15475.519999999999,"truck_hire":13900,"net_revenue":1575.5199999999986,"margins":0.10180724137217999,"total_miles":3616,"rev_per_mile":0.4357079646017695,"stop_percentage":96.875,"order_percentage":93.75},{"dispacher":"JOY LYNN","total_orders":153,"revenue":197783.47,"truck_hire":148254.8000000001,"net_revenue":49528.6699999999,"margins":0.25041865227665333,"total_miles":66821,"rev_per_mile":0.7412141392675939,"stop_percentage":86.75078864353313,"order_percentage":75.81699346405229},{"dispacher":"LENORA SMITH","total_orders":110,"revenue":184275.27000000002,"truck_hire":142752.5,"net_revenue":41522.77000000002,"margins":0.22533012704309163,"total_miles":62344,"rev_per_mile":0.6660267226998592,"stop_percentage":90.43062200956938,"order_percentage":81.81818181818183},{"dispacher":"LIZ SWENSON","total_orders":163,"revenue":372709.86,"truck_hire":306154.95,"net_revenue":66554.90999999997,"margins":0.17857029594011808,"total_miles":144309,"rev_per_mile":0.46119722262644725,"stop_percentage":86.01190476190477,"order_percentage":74.23312883435584},{"dispacher":"MIJKEN CASSIDY","total_orders":114,"revenue":171870.26,"truck_hire":139868.64999999994,"net_revenue":32001.610000000073,"margins":0.18619632041052403,"total_miles":62368,"rev_per_mile":0.513109447152387,"stop_percentage":95.78059071729957,"order_percentage":92.98245614035088},{"dispacher":"RIKI MAROTZ","total_orders":146,"revenue":233553.71,"truck_hire":200856.15000000002,"net_revenue":32697.55999999997,"margins":0.14000017383581692,"total_miles":84602,"rev_per_mile":0.3864868442826407,"stop_percentage":98.18181818181819,"order_percentage":96.57534246575342},{"dispacher":"SAM SWENSON","total_orders":149,"revenue":296703.33999999997,"truck_hire":247922.35,"net_revenue":48780.98999999996,"margins":0.16440997934165474,"total_miles":111282,"rev_per_mile":0.438354720439963,"stop_percentage":93.4931506849315,"order_percentage":87.24832214765101},{"dispacher":"Total","total_orders":1013,"revenue":1785742.6300000004,"truck_hire":1453362.1,"net_revenue":332380.53000000026,"margins":0.18613014239347592,"total_miles":652967,"rev_per_mile":0.509031130210256,"stop_percentage":90.16706443914082,"order_percentage":81.83613030602172}];


const locationGroups = {
    "Wellsvile": ["CAMI HANSEN", "LIZ SWENSON", "SAM SWENSON", "LENORA SMITH" ],
    "SLC": ["JOY LYNN", "MIJKEN CASSIDY"],
    "Ashton": ["JERRAMI MAROTZ", "RIKI MAROTZ"]

};

let groupByLocation = (data, groups) => {
	return Object.keys(groups).reduce((acc, location) => {
		acc[location] = data.filter(row => groups[location].includes(row.dispacher));
		return acc;
	}, {});
};

let groupedData = groupByLocation(data, locationGroups);

console.log(groupedData);



