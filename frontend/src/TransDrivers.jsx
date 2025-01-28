import React, { useState, useEffect } from 'react';
import {
  ComposedChart,
  Bar,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  CartesianGrid,
  ResponsiveContainer,
} from 'recharts';

const TransDrivers = () => {
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [drivers, setDrivers] = useState([]);
  const [fleetManagers, setFleetManagers] = useState([]);
  const [selectedDriver, setSelectedDriver] = useState('');
  const [selectedFleetManager, setSelectedFleetManager] = useState('All');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const apiURL = import.meta.env.VITE_API_URL; // Set your backend API URL

  // Fetch data from the updated API endpoint
  useEffect(() => {
    setLoading(true);
    setError(null);

    fetch(`${apiURL}/Transportation/DriverManager`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }
        return response.json();
      })
      .then((data) => {
        setData(data.Data);
        setDrivers([...new Set(data.Data.map((item) => item.driver_id))]); // Get unique drivers
        setFleetManagers([
          'All',
          ...new Set(data.Data.map((item) => item.fleet_manager.trim())),
        ]); // Include "All" option
        setFilteredData(data.Data); // Initially, show all data
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [apiURL]);

  // Filter logic
  useEffect(() => {
    let filtered = data;

    if (selectedFleetManager !== 'All') {
      filtered = filtered.filter(
        (item) => item.fleet_manager.trim() === selectedFleetManager
      );
    }

    if (selectedDriver) {
      filtered = filtered.filter((item) => item.driver_id === selectedDriver);
    }

    setFilteredData(filtered);
  }, [selectedFleetManager, selectedDriver, data]);

  // Linear regression function
  const calculateLinearRegression = (data) => {
    const n = data.length;
    const x = data.map((d, i) => i + 1); // X values: Week numbers (1, 2, 3, ...)
    const y = data.map((d) => d.totalDistance); // Y values: Total distances

    const xAvg = x.reduce((sum, xi) => sum + xi, 0) / n;
    const yAvg = y.reduce((sum, yi) => sum + yi, 0) / n;

    let numerator = 0;
    let denominator = 0;

    for (let i = 0; i < n; i++) {
      numerator += (x[i] - xAvg) * (y[i] - yAvg);
      denominator += (x[i] - xAvg) ** 2;
    }

    const slope = numerator / denominator;
    const intercept = yAvg - slope * xAvg;

    // Calculate the linear regression line
    const regressionLine = x.map((xi) => ({
      week: `Week ${xi}`,
      totalDistance: slope * xi + intercept,
    }));

    return regressionLine;
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  // Prepare chart data
  const chartData =
    filteredData.length > 0
      ? filteredData[0].summaries.map((total, index) => ({
        week: `Week ${index + 1}`,
        totalDistance: total,
      }))
      : [];

  const regressionData = calculateLinearRegression(chartData);

  return (
    <div className="p-8 bg-gray-100 min-h-screen">
      <h1 className="text-4xl font-extrabold mb-8 text-center">Driver Miles</h1>

      {/* Filters */}
      <div className="flex gap-4 mb-8">
        <select
          value={selectedFleetManager}
          onChange={(e) => {
            setSelectedFleetManager(e.target.value);
            setSelectedDriver(''); // Reset driver when fleet manager changes
          }}
          className="p-2 border rounded"
        >
          <option value="All">All Fleet Managers</option>
          {fleetManagers.map((manager) => (
            <option key={manager} value={manager}>
              {manager}
            </option>
          ))}
        </select>

        <select
          value={selectedDriver}
          onChange={(e) => setSelectedDriver(e.target.value)}
          className="p-2 border rounded"
        >
          <option value="">Select Driver</option>
          {drivers.map((driver) => (
            <option key={driver} value={driver}>
              {driver}
            </option>
          ))}
        </select>
      </div>

      {/* Chart */}
      {chartData.length > 0 ? (
        <ResponsiveContainer width="100%" height={400}>
          <ComposedChart data={chartData}>
            <CartesianGrid stroke="#f5f5f5" />
            <XAxis
              dataKey="week"
              label={{
                value: 'Week Number',
                position: 'insideBottom',
                offset: -5,
              }}
            />
            <YAxis
              label={{
                value: 'Distance (miles)',
                angle: -90,
                position: 'insideLeft',
              }}
            />
            <Tooltip
              formatter={(value, name) => [`${value} miles`, name]}
              labelFormatter={(label) => `Week: ${label}`}
            />
            <Legend />
            <Bar
              dataKey="totalDistance"
              barSize={20}
              fill="#8884d8"
              name="Total Distance"
            />
            <Line
              type="monotone"
              dataKey="totalDistance"
              data={regressionData}
              stroke="#ff7300"
              dot={false}
              name="Linear Regression"
            />
          </ComposedChart>
        </ResponsiveContainer>
      ) : (
        <p className="text-center text-gray-500">
          Select a driver or fleet manager to view the data.
        </p>
      )}
    </div>
  );
};

export default TransDrivers;
