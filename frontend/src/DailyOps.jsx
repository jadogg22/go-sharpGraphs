import React, { useState, useEffect } from 'react';
 import { PropagateLoader } from 'react-spinners';


const DailyOps = () => {

        const [data, setData] = useState([]);
        const [loading, setLoading] = useState(true);
        const [error, setError] = useState(null);
    
        useEffect(() => {
            setLoading(true);
            setError(null); // Reset error state before fetching
            fetch('http://192.168.0.62:5000/api/Transportation/Daily_Ops/')
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
            return <div className="text-center"><PropagateLoader /></div>;
        }
    
        if (error) {
            return <p className="text-red-500">Failed to load data: {error}</p>;
        }
    

      return (
        <div className="p-6 bg-grey-100 min-h-screen">
          <h1 className="text-3xl font-bold text-center mb-8">Daily Operations Summary</h1>
          <div className="overflow-x-auto">
            <table className="min-w-full bg-white border border-gray-200 rounded-lg shadow-sm">
              <thead className="bg-gray-800">
                <tr>
                  <th className="py-3 px-4 text-left text-sm font-semibold text-gray-100 uppercase tracking-wider">Driver Manager</th>
                  <th className="py-3 px-4 text-left text-sm font-semibold text-gray-100 uppercase tracking-wider">Number of Trucks</th>
                  <th className="py-3 px-4 text-left text-sm font-semibold text-gray-100 uppercase tracking-wider">Miles per Truck</th>
                  <th className="py-3 px-4 text-left text-sm font-semibold text-gray-100 uppercase tracking-wider">Deadhead</th>
                  <th className="py-3 px-4 text-left text-sm font-semibold text-gray-100 uppercase tracking-wider">Order</th>
                  <th className="py-3 px-4 text-left text-sm font-semibold text-gray-100 uppercase tracking-wider">Stop</th>
                </tr>
              </thead>
              <tbody>
                {data.map((row, index) => (
                  <tr key={index} className="border-t hover:bg-gray-100">
                    <td className="py-3 px-4 text-gray-900">{row.driverManager}</td>
                    <td className="py-3 px-4 text-gray-900">{row.numberOfTrucks}</td>
                    <td className="py-3 px-4 text-gray-900">{Math.round(row.milesPerTruck)}</td>
                    <td className="py-3 px-4 text-gray-900">{Math.round(row.deadhead)}%</td>
                    <td className="py-3 px-4 text-gray-900">{Math.round(row.order)}%</td>
                    <td className="py-3 px-4 text-gray-900">{Math.round(row.stop)}%</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      );
    };

export default DailyOps;
