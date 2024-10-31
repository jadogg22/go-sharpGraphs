
import React, { useState, useEffect } from 'react';
import {
  ResponsiveContainer,
  PieChart,
  Pie,
  Tooltip,
  Legend,
  Cell,
  BarChart,
  Bar,
  XAxis,
  YAxis,
} from 'recharts';

const dataLineChart = [
  { name: 'Jan', revenue: 4000 },
  { name: 'Feb', revenue: 3000 },
  { name: 'Mar', revenue: 5000 },
  { name: 'Apr', revenue: 4000 },
  { name: 'May', revenue: 6000 },
  { name: 'Jun', revenue: 7000 },
];

const dataBarChart = [
  { category: 'Product A', revenue: 2400 },
  { category: 'Product B', revenue: 3000 },
  { category: 'Product C', revenue: 2000 },
  { category: 'Product D', revenue: 2780 },
];

const aggregateHosStatus = (data) => {
  return data.reduce((acc, driver) => {
    const status = driver.HosStatus || 'Unknown';
    acc[status] = (acc[status] || 0) + 1;
    return acc;
  }, {});
};

const Dashboard = () => {
  // Initialize DriverData as an empty array
  const [DriverData, setDriverData] = useState([]);
  const [IsLoading, setIsLoading] = useState(false);
  const [Error, setError] = useState(null);

  const apiURL = import.meta.env.VITE_API_URL;

  const fetchData = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`${apiURL}/Transportation/dashboard`);
      if (!response.ok) {
        throw new Error(`Error: ${response.statusText}`);
      }
      const data = await response.json();
      setDriverData(data);
    } catch (error) {
      console.error('There was a problem with your fetch operation:', error);
      setError(error.toString());
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  // Aggregate only if DriverData is available
  const aggregatedData = aggregateHosStatus(DriverData);
  const chartData = Object.entries(aggregatedData).map(([key, value]) => ({
    name: key,
    value,
  }));

  // Define colors for each status
  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#FF6699'];

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <h1 className="text-3xl text-center text-gray-800 mb-16">Revenue Reporting Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div className="border-4 border-gray-600 bg-white p-4 shadow-lg rounded-lg">
          <h2 className="text-xl font-extrabold text-gray-700">Revenue</h2>
          <p className="text-2xl font-bold text-gray-800">$6,000</p>
          <p className="text-sm text-green-500">↑ 10% from last week</p>
        </div>
        <div className="border-4 border-gray-600 bg-white p-4 shadow-lg rounded-lg">
          <h2 className="text-lg text-gray-700">Orders Today</h2>
          <p className="text-2xl font-bold text-gray-800">150</p>
          <p className="text-sm text-red-500">↓ 5% from yesterday</p>
        </div>
        <div className="border-4 border-gray-600 bg-white p-4 shadow-lg rounded-lg">
          <h2 className="text-lg text-gray-700">Orders This Week</h2>
          <p className="text-2xl font-bold text-gray-800">800</p>
          <p className="text-sm text-green-500">↑ 15% from last week</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="border-4 border-gray-600 bg-white p-4 shadow-lg rounded-lg">
          <h2 className="text-lg font-bold text-gray-700 pb-4">Transportation Drivers' Status</h2>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={chartData}
                cx='50%'
                cy='50%'
                fill="#8884d8"
                paddingAngle={5}
                dataKey="value"
              >
                {chartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="border-4 border-gray-600 bg-white p-4 shadow-lg rounded-lg">
          <h2 className="text-lg font-bold text-gray-700 pb-4">Revenue by Product</h2>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={dataBarChart}>
              <XAxis dataKey="category" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="revenue" fill="#82ca9d" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;

