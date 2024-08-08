import React, { useState, useEffect } from 'react';


const DailyOps = () => {
  
    const [data] = useState([
        {
          driverManager: 'John Doe',
          numberOfTrucks: 12,
          milesPerTruck: 350,
          deadhead: 25,
          order: 'ORD12345',
          stop: 'STOP67890',
        },
        {
          driverManager: 'Jane Smith',
          numberOfTrucks: 8,
          milesPerTruck: 420,
          deadhead: 30,
          order: 'ORD23456',
          stop: 'STOP78901',
        },
        {
          driverManager: 'Mark Johnson',
          numberOfTrucks: 15,
          milesPerTruck: 300,
          deadhead: 20,
          order: 'ORD34567',
          stop: 'STOP89012',
        },
        {
          driverManager: 'Emily Davis',
          numberOfTrucks: 10,
          milesPerTruck: 375,
          deadhead: 15,
          order: 'ORD45678',
          stop: 'STOP90123',
        },
      ]);
    
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
                    <td className="py-3 px-4 text-gray-900">{row.milesPerTruck}</td>
                    <td className="py-3 px-4 text-gray-900">{row.deadhead}</td>
                    <td className="py-3 px-4 text-gray-900">{row.order}</td>
                    <td className="py-3 px-4 text-gray-900">{row.stop}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      );
    };

export default DailyOps;
