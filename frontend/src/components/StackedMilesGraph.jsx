import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const MilesBarChart = ({ data }) => {
  const formatTooltip = (value, name, props) => {
    return `${name}: ${value.toFixed(2)}`;
  };

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div style={{ backgroundColor: 'white', padding: '10px', border: '1px solid #ccc' }}>
          <p><strong>{label}</strong></p>
          <p>Total Loaded Miles: {data["Total_Loaded_Miles"].toFixed(2)}</p>
          <p>Total Empty Miles: {data["Total_Empty_Miles"].toFixed(2)}</p>
          <p>Total Miles: {data["Total_Actual_Miles"].toFixed(2)}</p>
          <p>Percent Empty: {data["Percent_empty"].toFixed(2)}</p>
        </div>
      );
    }
    return null;
  };

  return (
    <ResponsiveContainer width="100%" height={400}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="NameStr" />
        <YAxis label={{ value: 'Miles', angle: -90, position: 'insideLeft' }} />
        <Tooltip content={<CustomTooltip />} />
        <Legend />
        <Bar dataKey="TotalLoadedMiles" stackId="a" fill="#82ca9d" name="Loaded Miles" />
        <Bar dataKey="TotalEmptyMiles" stackId="a" fill="#ff8042" name="Empty Miles" />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default MilesBarChart;
