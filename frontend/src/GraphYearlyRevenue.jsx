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

  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');

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

  const downloadReport = () => {
    if (!startDate || !endDate) {
      setError('Please select both a start and end date.');
      return;
    }
    setLoading(true);
    setError(null);

    const url = `/api/Transportation/statistics/laneprofitability/report?startDate=${startDate}&endDate=${endDate}`;

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
      {/* Report Download Section */}
      <div className="flex items-center space-x-4 mb-4">
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
          onClick={downloadReport}
          disabled={isLoading}
          className="mt-5 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:bg-gray-400"
        >
          {isLoading ? 'Downloading...' : 'Download Lane Profitability Report'}
        </button>
      </div>
      {error && <p className="text-red-500">Error: {error}</p>}

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

