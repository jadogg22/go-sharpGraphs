import React, { useState } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, ReferenceLine, Cell, LabelList } from 'recharts';

const LaneProfitability = () => {
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [totalLoads, setTotalLoads] = useState(0);

  const fetchData = () => {
    if (!startDate || !endDate) {
      setError('Please select both a start and end date.');
      return;
    }
    setLoading(true);
    setError(null);

    const url = `/api/Transportation/LaneProfit?startDate=${startDate}&endDate=${endDate}`;

    fetch(url)
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(responseData => {
        setData(responseData.data);
        // Calculate total loads
        const loads = responseData.data.reduce((sum, item) => sum + item.total_trips, 0);
        setTotalLoads(loads);
        setLoading(false);
      })
      .catch(error => {
        setError(error.message);
        setLoading(false);
      });
  };

  const downloadReport = () => {
    if (!startDate || !endDate) {
      setError('Please select both a start and end date.');
      return;
    }
    setLoading(true);
    setError(null);

    const url = `/api/statistics/laneprofitability/report?startDate=${startDate}&endDate=${endDate}`;

    fetch(url)
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.blob();
      })
      .then(blob => {
        const url = window.URL.createObjectURL(new Blob([blob]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `lane_profitability_report_${startDate}_${endDate}.pdf`);
        document.body.appendChild(link);
        link.click();
        link.parentNode.removeChild(link);
        setLoading(false);
      })
      .catch(error => {
        setError(error.message);
        setLoading(false);
      });
  };

  return (
    <div className="p-4 pt-64">
      <h1 className="text-2xl font-bold mb-4">Lane Profitability Analysis</h1>
      <div className="flex items-center space-x-4 mb-4 mt-16">
        <div>
          <label htmlFor="startDate" className="block text-sm font-medium text-gray-700">Start Date</label>
          <input
            type="date"
            id="startDate"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            className="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          />
        </div>
        <div>
          <label htmlFor="endDate" className="block text-sm font-medium text-gray-700">End Date</label>
          <input
            type="date"
            id="endDate"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
            className="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          />
        </div>
        <button
          onClick={fetchData}
          disabled={loading}
          className="mt-5 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400"
        >
          {loading ? 'Loading...' : 'Analyze'}
        </button>
        <button
          onClick={downloadReport}
          disabled={loading}
          className="mt-5 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:bg-gray-400"
        >
          {loading ? 'Downloading...' : 'Download Report'}
        </button>
      </div>

      {error && <p className="text-red-500">Error: {error}</p>}
      {data.length > 0 && <h2 className="text-xl font-semibold mb-4">Total Loads: {totalLoads}</h2>}

      <div style={{ width: '100%', height: 500 }}>
        <ResponsiveContainer>
          <BarChart
            data={data}
            margin={{
              top: 20, right: 30, left: 20, bottom: 5,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="other_state" />
            <YAxis yAxisId="left" orientation="left" stroke="#8884d8" />
            <YAxis yAxisId="right" orientation="right" stroke="#82ca9d" tickFormatter={(value) => `${value.toFixed(2)}`} />
            <Tooltip formatter={(value, name, props) => {
              if (name === 'Avg. Revenue/Mile') return [`${value.toFixed(2)}`, name];
              if (name === 'Total Trips') return [value, name];
              return [value, name];
            }} labelFormatter={(label, payload) => {
              if (payload && payload.length > 0) {
                return `Lane: UT <-> ${label} | Total Trips: ${payload[0].payload.total_trips}`;
              }
              return `Lane: UT <-> ${label}`;
            }} />
            <Legend />
            <ReferenceLine y={2.79} yAxisId="right" stroke="red" strokeDasharray="3 3" label="Target" />
            <Bar yAxisId="right" dataKey="avg_round_trip_revenue" name="Avg. Revenue/Mile" >
              {data.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={getBarColor(entry.total_trips)} />
              ))}
              <LabelList dataKey="total_trips" position="top" />
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>

      <h2 className="text-xl font-bold mt-8 mb-4">Lane Quality Score</h2>
      <div style={{ width: '100%', height: 500 }}>
        <ResponsiveContainer>
          <BarChart
            data={data.slice().sort((a, b) => b.lane_quality_score - a.lane_quality_score)}
            margin={{
              top: 20, right: 30, left: 20, bottom: 5,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="other_state" />
            <YAxis />
            <Tooltip formatter={(value, name) => [value.toFixed(2), name]} />
            <Legend />
            <Bar dataKey="lane_quality_score" name="Lane Quality Score">
              {data.slice().sort((a, b) => b.lane_quality_score - a.lane_quality_score).map((entry, index) => (
                <Cell key={`cell-${index}`} fill={getScoreColor(entry.lane_quality_score, data)} />
              ))}
              <LabelList dataKey="lane_quality_score" position="top" formatter={(value) => value.toFixed(2)} />
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

// Helper function to determine bar color based on total_trips
const getBarColor = (totalTrips) => {
  if (totalTrips >= 1000) return '#004d00'; // Dark Green
  if (totalTrips >= 500) return '#008000';  // Green
  if (totalTrips >= 200) return '#32cd32';  // Lime Green
  if (totalTrips >= 100) return '#adff2f';  // Green Yellow
  if (totalTrips >= 50) return '#ffd700';   // Gold
  if (totalTrips >= 20) return '#ffa500';   // Orange
  if (totalTrips >= 10) return '#ff4500';   // Orange Red
  return '#ff0000'; // Red
};

// Helper function to determine bar color based on LaneQualityScore (blue to yellow gradient)
const getScoreColor = (score, allData) => {
  // Ensure allData is not empty and contains valid numbers
  const validScores = allData.map(d => d.lane_quality_score).filter(s => typeof s === 'number' && !isNaN(s));

  if (validScores.length === 0) {
    return 'rgb(128, 128, 128)'; // Return a neutral color if no valid scores
  }

  const minScore = Math.min(...validScores);
  const maxScore = Math.max(...validScores);

  // Handle case where all scores are the same to avoid division by zero
  if (maxScore === minScore) {
    return 'rgb(0, 0, 139)'; // Return dark blue if all scores are identical
  }

  // Normalize score to a 0-1 range
  const normalizedScore = (score - minScore) / (maxScore - minScore);

  // Interpolate between blue and yellow
  const blue = [0, 0, 139]; // Dark Blue
  const yellow = [255, 255, 0]; // Yellow

  const r = Math.round(blue[0] + (yellow[0] - blue[0]) * normalizedScore);
  const g = Math.round(blue[1] + (yellow[1] - blue[1]) * normalizedScore);
  const b = Math.round(blue[2] + (yellow[2] - blue[2]) * normalizedScore);

  return `rgb(${r},${g},${b})`;
};

export default LaneProfitability;
