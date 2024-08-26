import React, { useState, useEffect } from 'react';
import { PropagateLoader } from 'react-spinners';

// Utility function to format numbers
const formatNumber = (number) => {
    if (number === null || number === undefined) return '';
    return new Intl.NumberFormat().format(number);
};

// Define dispatcher groups by location
const locationGroups = {
    "Wellsvile": ["CAMI HANSEN", "LIZ SWENSON", "SAM SWENSON", "LENORA SMITH"],
    "SLC": ["JOY LYNN", "MIJKEN CASSIDY"],
    "Ashton": ["JERRAMI MAROTZ", "RIKI MAROTZ"]
};

// Function to group and aggregate data by location
const groupAndAggregateByLocation = (data, groups) => {
    const locationData = Object.keys(groups).reduce((acc, location) => {
        const locationDispatchers = groups[location];
        const dispatchersData = data.filter(row => locationDispatchers.includes(row.dispacher));

        // Calculate totals for each location
        const totals = dispatchersData.reduce((sum, row) => {
            sum.total_orders += row.total_orders;
            sum.revenue += row.revenue;
            sum.truck_hire += row.truck_hire;
            sum.net_revenue += row.net_revenue;
            sum.margins += row.margins * row.total_orders; // Weighted average
            sum.total_miles += row.total_miles;
            sum.rev_per_mile += row.rev_per_mile * row.total_miles; // Weighted average
            sum.stop_percentage += row.stop_percentage * row.total_orders; // Weighted average
            sum.order_percentage += row.order_percentage * row.total_orders; // Weighted average
            return sum;
        }, {
            total_orders: 0,
            revenue: 0,
            truck_hire: 0,
            net_revenue: 0,
            margins: 0,
            total_miles: 0,
            rev_per_mile: 0,
            stop_percentage: 0,
            order_percentage: 0
        });

        // Calculate weighted averages
        const totalOrders = dispatchersData.reduce((sum, row) => sum + row.total_orders, 0);
        if (totalOrders > 0) {
            totals.margins /= totalOrders;
            totals.rev_per_mile /= totals.total_miles;
            totals.stop_percentage /= totalOrders;
            totals.order_percentage /= totalOrders;
        }

        acc[location] = {
            dispatchers: dispatchersData,
            totals: totals
        };
        return acc;
    }, {});

    return locationData;
};

const LogisticsMTD = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [groupedData, setGroupedData] = useState({});

    useEffect(() => {
        setLoading(true);
        setError(null); // Reset error state before fetching
        fetch('http://192.168.0.62:5000/api/Logistics/MTD')
            .then(response => {
                if (!response.ok) {
                    throw new Error(`Error: ${response.statusText}`);
                }
                return response.json();
            })
            .then(data => {
                const locationData = groupAndAggregateByLocation(data, locationGroups);
                setData(data);
                setGroupedData(locationData);
                setLoading(false);
            })
            .catch(err => {
                setError(err.message);
                setLoading(false);
            });
    }, []);

    if (loading) {
        return <div className="flex justify-center items-center min-h-screen"><PropagateLoader /></div>;
    }

    if (error) {
        return <p className="text-center text-red-500 text-lg">Failed to load data: {error}</p>;
    }

    return (
        <div className="p-8 bg-gray-100 min-h-screen flex flex-col items-center">
            <h1 className="text-4xl font-extrabold mb-8">Logistics Month to Date Stats</h1>
            {Object.keys(groupedData).map(location => (
                <div key={location} className="w-full mb-8">
                    <h2 className="text-2xl font-bold mb-4">{location}</h2>
                    <div className="overflow-x-auto">
                        <table className="min-w-full bg-white border border-gray-200 rounded-lg shadow-lg">
                            <thead className="bg-gray-800 text-white">
                                <tr>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Dispatcher</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Total Orders</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Revenue</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Truck Hire</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Net Revenue</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Margins (%)</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Total Miles</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Rev Per Mile</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Stop Percentage (%)</th>
                                    <th className="py-4 px-6 text-left text-lg font-semibold">Order Percentage (%)</th>
                                </tr>
                            </thead>
                            <tbody>
                                {groupedData[location].dispatchers.map((row, index) => (
                                    <tr key={index} className="border-t hover:bg-gray-100">
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{row.dispacher}</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{formatNumber(row.total_orders)}</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{formatNumber(row.revenue.toFixed(2))}</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{formatNumber(row.truck_hire.toFixed(2))}</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{formatNumber(row.net_revenue.toFixed(2))}</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{(row.margins * 100).toFixed(2)}%</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{formatNumber(row.total_miles)}</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{formatNumber(row.rev_per_mile.toFixed(2))}</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{row.stop_percentage.toFixed(2)}%</td>
                                        <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{row.order_percentage.toFixed(2)}%</td>
                                    </tr>
                                ))}
                                {/* Totals row */}
                                <tr className="border-t bg-gray-100 font-bold">
                                    <td className="py-4 px-6 text-gray-900 text-lg">Totals</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{formatNumber(groupedData[location].totals.total_orders)}</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{formatNumber(groupedData[location].totals.revenue.toFixed(2))}</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{formatNumber(groupedData[location].totals.truck_hire.toFixed(2))}</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{formatNumber(groupedData[location].totals.net_revenue.toFixed(2))}</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{(groupedData[location].totals.margins * 100).toFixed(2)}%</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{formatNumber(groupedData[location].totals.total_miles)}</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{formatNumber(groupedData[location].totals.rev_per_mile.toFixed(2))}</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{groupedData[location].totals.stop_percentage.toFixed(2)}%</td>
                                    <td className="py-4 px-6 text-gray-900 text-lg">{groupedData[location].totals.order_percentage.toFixed(2)}%</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            ))}
        </div>
    );
};

export default LogisticsMTD;