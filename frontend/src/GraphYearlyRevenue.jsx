import {
  ResponsiveContainer,
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  Tooltip,
  Legend,
  YAxis
} from 'recharts';
import React, { useState, useEffect } from 'react';
import { Menu } from '@headlessui/react';
import { PropagateLoader } from 'react-spinners';

function formatNumberWithCommas(number) {
  number = number.toFixed(2);
  return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

const CustomTooltip = ({ active, payload, label, selectedYears, colors, years }) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;

    return (
      <div style={{ background: 'white', padding: '5px', border: '1px solid #ccc' }}>
        <p>{`Week: ${data.Name}`}</p>
        {selectedYears.map((year, index) => {
          const revenueKey = `${year} Revenue`;
          if (data[revenueKey] !== undefined && data[revenueKey] !== null) {
            const colorIndex = years.indexOf(year);
            return (
              <p key={year} style={{ color: colors[colorIndex] }}>
                {`Revenue ${year}: $${formatNumberWithCommas(data[revenueKey])}`}
              </p>
            );
          }
          return null;
        })}
      </div>
    );
  }

  return null;
};

const YearlyRevenue = ({ company }) => {
  const currentYear = new Date().getFullYear();
  const years = Array.from({ length: currentYear - 2019 }, (_, i) => (2020 + i).toString());
  const colors = ["#ef4444", "#fbbf24", "#0ea5e9", "#50c878", "#8b5cf6", "#ec4899", "#14b8a6", "#f97316", "#6b7280", "#ef4444"];

  const [data, setData] = useState(null);
  const [isLoading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedYears, setSelectedYears] = useState(years);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const response = await fetch(`/api/Transportation/get_yearly_revenue`);
        if (!response.ok) {
          setError('Network response was not ok');
          throw new Error('Network response was not ok');
        }
        const result = await response.json();
        const transformedData = result.Data.map(item => {
          const newItem = { Name: item.Name };
          for (const key in item.Revenues) {
            newItem[key] = item.Revenues[key];
          }
          return newItem;
        });
        setData(transformedData);
      } catch (error) {
        console.error('There was a problem with your fetch operation:', error);
        setError(error.toString());
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleYearChange = (year) => {
    setSelectedYears((prevSelectedYears) =>
      prevSelectedYears.includes(year)
        ? prevSelectedYears.filter((y) => y !== year)
        : [...prevSelectedYears, year]
    );
  };

  if (isLoading) {
    return <div className="text-center"><PropagateLoader /></div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!data) {
    return <div>No data available.</div>;
  }

  const now = new Date();
  const startOfYear = new Date(now.getFullYear(), 0, 1);
  const currentWeek = Math.ceil((((now - startOfYear) / 86400000) + startOfYear.getDay() + 1) / 7);
  const currentYearString = now.getFullYear().toString();

  const chartData = data.map(d => {
    const newD = { ...d };
    const itemWeek = parseInt(newD.Name, 10);
    const revenueKey = `${currentYearString} Revenue`;

    if (selectedYears.includes(currentYearString) && itemWeek > currentWeek) {
      newD[revenueKey] = null;
    }
    return newD;
  });


  return (
    <div className="w-full h-[500px]">
      {/* Year Filter Menu */}
      <Menu as="div" className="mb-4">
        <div className="flex flex-wrap gap-2">
          {years.map((year) => (
            <button
              key={year}
              onClick={() => handleYearChange(year)}
              className={`px-3 py-1 border rounded-full ${selectedYears.includes(year) ? 'bg-blue-600 text-white' : 'bg-gray-200'}`}
            >
              {year}
            </button>
          ))}
        </div>
      </Menu>

      {/* Chart */}
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={chartData}>
          <CartesianGrid stroke="#ccc" />
          <XAxis dataKey="Name" label={{ value: 'Week', position: 'insideBottomRight', offset: -5 }} />
          <YAxis tickFormatter={(value) => `${value / 1000}k`} />
          <Tooltip
            content={
              <CustomTooltip
                selectedYears={selectedYears}
                colors={colors}
                years={years}
              />
            }
          />
          <Legend />
          {selectedYears.map((year, index) => (
            <Line
              key={year}
              type="monotone"
              dataKey={`${year} Revenue`}
              stroke={colors[years.indexOf(year) % colors.length]}
              dot={false}
              activeDot={{ r: 6 }} // 👈 Enlarges point on hover
              strokeWidth={3}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default YearlyRevenue;

