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
    //first only grab the number to the hundredth place
    number = number.toFixed(2);
    return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  }
  
  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      const data = payload.map((entry) => entry.payload);
  
      return (
        <div style={{ background: 'white', padding: '5px', border: '1px solid #ccc' }}>
          <p>{`Week: ${data[0].Name}`}</p>
          {data[0]["2021 Revenue"] !== undefined && (
            <p className='text-[#ef4444]'>{`Revenue 2021: $ ${formatNumberWithCommas(data[0]["2021 Revenue"])}`}</p>
          )}
          {data[0]["2022 Revenue"] !== undefined && (
            <p className='text-[#fbbf24]'>{`Revenue 2022: $ ${formatNumberWithCommas(data[0]["2022 Revenue"])}`}</p>
          )}
          {data[0]["2023 Revenue"] !== undefined && (
            <p className='text-[#0ea5e9]'>{`Revenue 2023: $ ${formatNumberWithCommas(data[0]["2023 Revenue"])}`}</p>
          )}
          {data[0]["2024 Revenue"] !== undefined && (
            <p className='text-[#50c878]'>{`Revenue 2024: $ ${formatNumberWithCommas(data[0]["2024 Revenue"])}`}</p>
          )}
        </div>
      );
    }
  
    return null;
  };
  
  const YearlyRevenue = ({ company }) => {
    const years = ["2021", "2022", "2023", "2024"];
    const colors = ["#ef4444", "#fbbf24", "#0ea5e9", "#50c878"];
  
    const [data, setData] = useState(null);
    const [isLoading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [selectedYears, setSelectedYears] = useState(years); // State for selected years
  
    const convertData = (data) => {
      return data.map((entry) => {
        return {
          ...entry,
          '2021 Revenue': parseInt(entry['2021 Revenue']),
          '2022 Revenue': parseInt(entry['2022 Revenue']),
          '2023 Revenue': parseInt(entry['2023 Revenue']),
          '2024 Revenue': parseInt(entry['2024 Revenue']),
        };
      });
    };
  
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const response = await fetch('http://192.168.0.62:5000/api/Transportation/get_yearly_revenue');
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const result = await response.json();
        setData(result.Data);
      } catch (error) {
        console.error('There was a problem with your fetch operation:', error);
        setError(error.toString());
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);  
    if (isLoading) {
      return <div className="text-center"><PropagateLoader /></div>;
    }
  
    if (error) {
      return <div>Error: {error.message}</div>;
    }
  
    if (!data) {
      return <div>No data available.</div>;
    }
  
    const handleYearChange = (year) => {
      setSelectedYears((prevSelectedYears) => {
        if (prevSelectedYears.includes(year)) {
          return prevSelectedYears.filter((y) => y !== year);
        } else {
          return [...prevSelectedYears, year];
        }
      });
    };
  
    return (
      <div className="relative">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-800 pb-4">Transportation Weekly Revenue (in Thousands)</h1>
        </div>
  
        <div className="text-right mb-2 mr-4 z-50">
          <Menu as="div" className="relative inline-block text-left">
            <div>
              <Menu.Button className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
                Select Years
              </Menu.Button>
            </div>
            <Menu.Items className="origin-top-right absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-50">
              {years.map((year) => (
                <Menu.Item key={year}>
                  {({ active }) => (
                    <div
                      className={`${
                        active ? 'bg-gray-100' : ''
                      } flex items-center px-4 py-2 text-gray-700 cursor-pointer`}
                      onClick={(event) => event.stopPropagation()} // Prevent the default action
                    >
                      <input
                        type="checkbox"
                        value={year}
                        checked={selectedYears.includes(year)}
                        onChange={() => handleYearChange(year)}
                        className="mr-2"
                      />
                      {year}
                    </div>
                  )}
                </Menu.Item>
              ))}
            </Menu.Items>
          </Menu>
        </div>
  
        <ResponsiveContainer className="bg-gray-50" width="100%" height={450}>
          <LineChart data={data} margin={{ top: 20, left: 20, right: 20, bottom: 50 }}>
            <CartesianGrid strokeDasharray="3" stroke="#ccc" />
            <XAxis dataKey="Name" stroke="black" />
            <YAxis
              stroke="black"
              domain={([dataMin, dataMax]) => {
                const roundedMin = Math.floor(dataMin / 100000) * 100000;
                const roundedMax = Math.ceil(dataMax / 50000) * 50000;
                return [roundedMin, roundedMax];
              }}
              tickFormatter={(value) => `${Math.round(value / 1000)}K`}
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
  
            {selectedYears.map((year, index) => (
              <Line
                key={year}
                type="monotone"
                dataKey={`${year} Revenue`}
                name={`${year} Revenue`}
                stroke={colors[index]}
                strokeWidth={3}
              />
            ))}
          </LineChart>
        </ResponsiveContainer>
      </div>
    );
  };
  
  export default YearlyRevenue;
