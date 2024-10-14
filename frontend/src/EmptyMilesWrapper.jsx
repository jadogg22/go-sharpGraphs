import React, { useState, useEffect } from 'react';
import StackedMilesGraph from './components/StackedMilesGraph';


const EmptyMilesWrapper = () => {
  const [milesData, setMilesData] = useState({
    week_to_date: [],
    month_to_date: []
  });

  const [timeFrame, setTimeFrame] = useState('week_to_date');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const apiURL = import.meta.env.VITE_API_URL;

  const fetchData = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`${apiURL}/Transportation/Stacked_miles/${timeFrame}`);
      if (!response.ok) {
        console.error(`Error: ${response.statusText}`);
        setError(`Error: ${response.statusText}`);
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      if (timeFrame === 'month_to_date') {
        setMilesData((prevData) => ({
          ...prevData,
          month_to_date: data.month_to_date || [],
        }));
      } else {
        setMilesData((prevData) => ({
          ...prevData,
          week_to_date: data.week_to_date || [],
        }));
      }

    } catch (error) {
      console.error('There was a problem with your fetch operation:', error);
      setError(error.toString());
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [timeFrame]);

  const toggleTimeFrame = () => {
    setTimeFrame(prevTimeFrame =>
      prevTimeFrame === 'month_to_date' ? 'week_to_date' : 'month_to_date'
    );
  };

  if (isLoading) {
    return <p>Loading...</p>;
  }

  if (error) {
    return <p style={{ color: 'red' }}>{error}</p>;
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-center p-6">Empty Miles - {timeFrame === 'month_to_date' ? "Weekly View" : "Daily view"}</h1>
      <div className="flex justify-between items-center px-16 pb-6"> {/* Flex container for button */}
        <button
          onClick={toggleTimeFrame}
          className="bg-blue-400 text-white font-bold py-2 px-4 rounded focus:outline-none"
          style={{ alignSelf: 'flex-end' }}  // Align button to the bottom of its container
        >
          Switch to {timeFrame === 'month_to_date' ? 'Week to Date' : 'Month to Date'}
        </button>
      </div>
      {isLoading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {!isLoading && !error && timeFrame === 'month_to_date' &&
        <StackedMilesGraph data={milesData.month_to_date} />}
      {!isLoading && !error && timeFrame === 'week_to_date' &&
        <StackedMilesGraph data={milesData.week_to_date} />}
    </div>
  );
};

export default EmptyMilesWrapper;
