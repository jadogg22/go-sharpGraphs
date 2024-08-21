import React, { useState, useEffect } from 'react';
import { PropagateLoader } from 'react-spinners';



const DailyOps = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const fetchData = async () => {
        setLoading(true);
        setError(null); // Reset error state before fetching
        fetch('http://192.168.0.62:5000/api/Transportation/Daily_Ops')
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
        };

        
    useEffect(() => {
        fetchData();

        const interval = setInterval(() => {
            fetchData();
        }, 45 * 60 * 1000);

        return () => clearInterval(interval);
    }, []);

    if (loading) {
        return <div className="flex justify-center items-center min-h-screen"><PropagateLoader /></div>;
    }

    if (error) {
        return <p className="text-center text-red-500 text-lg">Failed to load data: {error}</p>;
    } 

    return (
        <div className="p-6 bg-gray-100 min-h-screen flex flex-col items-center">
            <h1 className="text-4xl font-extrabold mb-8">Daily Operations Summary</h1>
            <div className="w-full max-w-6xl overflow-x-auto">
                <table className="min-w-full bg-white border border-gray-200 rounded-lg shadow-lg">
                    <thead className="bg-gray-800 text-white">
                        <tr>
                            <th className="py-4 px-6 text-left text-lg font-semibold">Driver Manager</th>
                            <th className="py-4 px-6 text-left text-lg font-semibold">Number of Trucks</th>
                            <th className="py-4 px-6 text-left text-lg font-semibold">Miles per Truck</th>
                            <th className="py-4 px-6 text-left text-lg font-semibold">Deadhead</th>
                            <th className="py-4 px-6 text-left text-lg font-semibold">Order</th>
                            <th className="py-4 px-6 text-left text-lg font-semibold">Stop</th>
                        </tr>
                    </thead>
                    <tbody>
                        {data.map((row, index) => (
                            <tr key={index} className="border-t hover:bg-gray-100">
                                <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{row.driverManager}</td>
                                <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{row.numberOfTrucks}</td>
                                <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{Math.round(row.milesPerTruck)}</td>
                                <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{Math.round(row.deadhead)}%</td>
                                <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{Math.round(row.order)}%</td>
                                <td className="py-4 px-6 text-gray-900 text-lg font-semibold">{Math.round(row.stop)}%</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default DailyOps;
