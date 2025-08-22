import React from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
import { PropagateLoader } from 'react-spinners';
import { useState, useEffect } from 'react';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

const RevPie = () => {
  const [codeData, setData] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div style={{ backgroundColor: 'white', padding: '10px', border: '1px solid #ccc' }}>
          <p><strong>{label}</strong></p>
          <p>Code: <strong>{data["Name"]}</strong></p>
          <p>Revenue: <strong>{data["Revenue"].toFixed(2)}</strong></p>
          <p>Number of Loads: <strong>{data["Count"]}</strong></p>
        </div>
      );
    }
    return null;
  };

  const CustomLabel = ({ index }) => {
    var name = codeData[index]["Name"];
    if (name.length > 15) {
      name = name.substring(0, 15) + "...";
    }
    var revenue = codeData[index]["Revenue"];
    return `${name} ${revenue.toFixed(2)}`;
  };

  const fetchData = async () => {
    setIsLoading(true);
    fetch(`/api/Transportation/get_coded_revenue/month`)
      .then(response => {
        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }
        return response.json();
      })
      .then(data => {
        setData(data["data"]);
        setIsLoading(false);
      })
      .catch(err => {
        setError(err.message);
        setIsLoading(false);
      });
  };

  useEffect(() => {
    fetchData();
  }, []);

  if (isLoading) {
    return <div className="text-center"><PropagateLoader /></div>;
  }
  if (error) {
    return <div>Error: {error}</div>;
  }
  if (!codeData) {
    return <div>No data available. Refresh?</div>;
  }
  return (
    <div className='p-8'>
      <p className='text-center text-2xl pb-8'>Revenue Pie Chart</p>
      <ResponsiveContainer width='100%' height={650}>
        <PieChart>
          <Pie
            data={codeData}
            dataKey='Revenue'
            nameKey='Name'
            cx='50%'
            cy='50%'
            //outerRadius={300}
            fill='#8884d8'
            labelLine={false}
            label={CustomLabel}
            animationDuration={500}
          >
            {
              codeData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))
            }
          </Pie>
          <Tooltip content={<CustomTooltip />} />
        </PieChart>
      </ResponsiveContainer>
    </div >
  );
}


export default RevPie;
