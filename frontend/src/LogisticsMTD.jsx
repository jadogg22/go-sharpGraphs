import React, { useState, useEffect } from 'react';
import { PropagateLoader } from 'react-spinners';

const formatNumber = (number) => {
    if (number === null || number === undefined) return '';
    return new Intl.NumberFormat().format(number);
};

const LogisticsMTD = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

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
                setData(data);
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
            <h1 className="text-4xl font-extrabold mb-8">Logistics Month to date stats</h1>
            <div className="w-full overflow-x-auto">
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
                        {data.map((row, index) => (
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
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default LogisticsMTD;
