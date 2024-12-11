import React, { useState } from "react";

const OrdersPage = () => {
  const [orderNumber, setOrderNumber] = useState("");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const apiURL = import.meta.env.VITE_API_URL;

  const downloadFile = async () => {
    setLoading(true);
    setError(""); // Reset error before trying to fetch
    try {
      // Ensure startDate and endDate are in the correct format (ISO 8601 format).
      const formattedStartDate = new Date(startDate).toISOString();
      const formattedEndDate = new Date(endDate).toISOString();

      const response = await fetch(`${apiURL}/Transportation/Generate_Sportsmans`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          orderNumber: orderNumber,
          startDate: formattedStartDate,
          endDate: formattedEndDate,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to generate Excel file");
      }

      const blob = await response.blob();
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = "Sportsmans_" + orderNumber + ".xlsx";
      link.click();
    } catch (error) {
      console.error("Error downloading file:", error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-50">
      <div className="bg-white p-8 rounded-lg shadow-lg w-96">
        <h2 className="text-xl font-semibold mb-4">Sportsmans Invoice</h2>

        <div className="mb-4">
          <label htmlFor="orderNumber" className="block text-sm font-medium text-gray-700">
            Order Number
          </label>
          <input
            type="text"
            id="orderNumber"
            value={orderNumber}
            onChange={(e) => setOrderNumber(e.target.value)}
            className="mt-2 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div className="mb-4">
          <label htmlFor="startDate" className="block text-sm font-medium text-gray-700">
            Start-Date
          </label>
          <input
            type="date"
            id="startDate"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            className="mt-2 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div className="mb-6">
          <label htmlFor="endDate" className="block text-sm font-medium text-gray-700">
            End-Date
          </label>
          <input
            type="date"
            id="endDate"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
            className="mt-2 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <button
          onClick={downloadFile}
          disabled={loading || !startDate || !endDate}
          className={`w-full py-2 text-white font-semibold rounded-md shadow-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 ${loading || !startDate || !endDate
            ? "bg-gray-400 cursor-not-allowed"
            : "bg-indigo-600 hover:bg-indigo-700"
            }`}
        >
          {loading ? "Generating..." : "Generate Invoice"}
        </button>

        {/* Error Popup Modal */}
        {error && (
          <div className="fixed inset-0 flex items-center justify-center bg-gray-500 bg-opacity-50 z-50">
            <div className="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
              <h3 className="text-lg font-semibold text-red-600">Error</h3>
              <p className="mt-2 text-sm text-gray-700">{error}</p>
              <button
                onClick={() => setError("")}
                className="mt-4 w-full py-2 text-white bg-red-500 hover:bg-red-600 rounded-md"
              >
                Close
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default OrdersPage;

