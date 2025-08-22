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

  // Fetch data from the updated API endpoint
  useEffect(() => {
    setLoading(true);
    setError(null);

    fetch(`/api/Transportation/DriverManager`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }
        return response.json();
      })
      .then((data) => {
        const sortedData = data.Data.sort((a, b) => a.driver_id.localeCompare(b.driver_id));
        setData(sortedData);
        setDrivers([...new Set(sortedData.map((item) => item.driver_id))]); // Get unique drivers
        const uniqueFleetManagers = [...new Set(data.Data.map((item) => item.fleet_manager.trim()))];
        setFleetManagers(['All', ...uniqueFleetManagers.sort()]); // Include "All" option and sort
        setFilteredData(sortedData); // Initially, show all data
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  // Filter logic
  useEffect(() => {
    let filtered = data;

    if (selectedFleetManager !== 'All') {
      filtered = filtered.filter(
        (item) => item.fleet_manager.trim() === selectedFleetManager
      );
    }

    // Update the drivers dropdown based on the selected fleet manager
    const availableDrivers = selectedFleetManager === 'All'
      ? [...new Set(data.map((item) => item.driver_id))]
      : [...new Set(filtered.map((item) => item.driver_id))];
    setDrivers(availableDrivers.sort());


    if (selectedDriver) {
      filtered = filtered.filter((item) => item.driver_id === selectedDriver);
    }

    setFilteredData(filtered);
  }, [selectedFleetManager, selectedDriver, data]);

  // Linear regression function
  const calculateLinearRegression = (data) => {
    const firstNonZeroIndex = data.findIndex(item => item.totalDistance > 0);

    if (firstNonZeroIndex === -1) {
      return [];
    }

    const regressionPoints = data.slice(firstNonZeroIndex);
    const n = regressionPoints.length;

    if (n < 2) {
      return [];
    }

    const x = Array.from(Array(n).keys()); // x will be [0, 1, 2, ...]
    const y = regressionPoints.map((d) => d.totalDistance);

    const sumX = x.reduce((a, b) => a + b, 0);
    const sumY = y.reduce((a, b) => a + b, 0);
    const sumXY = x.reduce((sum, xi, i) => sum + xi * y[i], 0);
    const sumXX = x.reduce((sum, xi) => sum + xi * xi, 0);

    const slope = (n * sumXY - sumX * sumY) / (n * sumXX - sumX * sumX);
    const intercept = (sumY - slope * sumX) / n;

    const regressionLine = data.map((d, i) => {
      let value = 0;
      if (i >= firstNonZeroIndex) {
        value = slope * (i - firstNonZeroIndex) + intercept;
      }
      return {
        week: `Week ${i + 1}`,
        totalDistance: Math.max(0, value),
      };
    });

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
              formatter={(value, name) => [`${Math.round(value)} miles`, name]}
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
