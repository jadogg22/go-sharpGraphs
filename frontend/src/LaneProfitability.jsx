import React, { useState } from 'react';
import { ClipLoader } from 'react-spinners';

const LaneProfitability = () => {
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

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

  return (
    <div className="flex flex-col items-center">
      <h1 className="text-2xl font-bold mb-4">Lane Profitability Report</h1>
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
          disabled={loading}
          className="mt-5 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:bg-gray-400 flex items-center"
        >
          {loading ? (
            <>
              <ClipLoader color={"#fff"} size={20} />
              <span className="ml-2">Downloading...</span>
            </>
          ) : (
            'Download Report'
          )}
        </button>
      </div>

      {error && <p className="text-red-500">Error: {error}</p>}
      <p className="text-gray-600 mt-4">Select a date range and click "Download Report" to generate the Lane Profitability PDF.</p>
    </div>
  );
};

export default LaneProfitability;
