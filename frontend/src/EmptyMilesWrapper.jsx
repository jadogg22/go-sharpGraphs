import React, { useState, useEffect } from 'react';
import axios from 'axios';
import StackedMilesGraph from './components/StackedMilesGraph';


const EmptyMilesWrapper = () => {
  const [milesData, setMilesData] = useState([]);
  const [timeFrame, setTimeFrame] = useState('month_to_date');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchData = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await axios.get(`http://192.168.0.62:5000/api/Transportation/Stacked_miles/${timeFrame}`);
      setMilesData(response.data);
      console.log('Data fetched:', response.data);
    } catch (err) {
      setError('Failed to fetch data. Please try again.');
      console.error('Error fetching data:', err);
    }
    setIsLoading(false);
  };

  useEffect(() => {
    fetchData();
  }, [timeFrame]);

  const toggleTimeFrame = () => {
    setTimeFrame(prevTimeFrame => 
      prevTimeFrame === 'month_to_date' ? 'week_to_date' : 'month_to_date'
    );
  };

  return (
    <div>
      <h1>Empty Miles</h1>
      <button onClick={toggleTimeFrame}>
        Switch to {timeFrame === 'month_to_date' ? 'Week to Date' : 'Month to Date'}
      </button>
      {isLoading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {!isLoading && !error && milesData.length > 0 && (
        <StackedMilesGraph data={milesData} />
      )}
    </div>
  );
};

export default EmptyMilesWrapper;
