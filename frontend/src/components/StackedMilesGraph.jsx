import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { useEffect } from 'react';

const MilesBarChart = ({ data }) => {


  const formatData = (data) => {
    data.forEach((entry) => {
      var date = entry["DeliveryDate"];
      entry["NameStr"] = formatDateShort(date);
      entry["DeliveryDate"] = formatDate(date);
    });
  };

  const formatTooltip = (value, name, props) => {
    return `${name}: ${value.toFixed(2)} `;
  };
  const formatDateShort = (date) => {
    var dateOBJ = new Date(date);

    if (dateOBJ === null) {
      return date;
    }
    // add a day to the date to fix the timezone issue
    dateOBJ.setDate(dateOBJ.getDate() + 1);
    const options = { month: 'short', day: 'numeric' };
    return dateOBJ.toLocaleDateString('en-US', options);
  };
  const formatDate = (date) => {
    var dateOBJ = new Date(date);
    if (dateOBJ === null) {
      return date;
    }

    dateOBJ.setDate(dateOBJ.getDate() + 1);
    const options = {
      weekday: 'short', // "Sun", "Mon", etc.
      year: 'numeric', // "2024"
      month: 'long', // "October"
      day: 'numeric' // "6"
    };
    return dateOBJ.toDateString();
  };


  useEffect(() => {
    formatData(data);
  }, [data]);

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div style={{ backgroundColor: 'white', padding: '10px', border: '1px solid #ccc' }}>
          <p><strong>{label}</strong></p>
          <p>Delevery Date: {data["DeliveryDate"]} </p>
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
    <ResponsiveContainer width="100%" height={600}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="NameStr" />
        <YAxis label={{ value: 'Miles', angle: -90, position: 'insideLeft' }} />
        <Tooltip content={<CustomTooltip />} />
        <Legend />
        <Bar dataKey="Total_Loaded_Miles" stackId="a" fill="#8edf7C" name="Loaded Miles" />
        <Bar dataKey="Total_Empty_Miles" stackId="a" fill="#FF0A12" name="Empty Miles" />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default MilesBarChart;
