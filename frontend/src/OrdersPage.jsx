import React, { useState } from "react";

const OrdersPage = () => {
  const [orderNumber, setOrderNumber] = useState("");
  const [orderDate, setOrderDate] = useState("");
  const [deliveryDate, setDeliveryDate] = useState("");
  const [loading, setLoading] = useState(false);

  const apiURL = import.meta.env.VITE_API_URL;



  const downloadFile = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${apiURL}/Transportation/Generate_Sportsmans`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          OrderNumber: orderNumber,
          StartDate: orderDate,
          EndDate: deliveryDate,
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
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-50">
      <div className="bg-white p-8 rounded-lg shadow-lg w-96">
        <h2 className="text-xl font-semibold mb-4">Sportsmans Invoice</h2>

        <div className="mb-4">
          <label htmlFor="deliveryDate" className="block text-sm font-medium text-gray-700">
            Order Number
          </label>
          <input
            type="text"
            id="OrderNumber"
            value={orderNumber}
            onChange={(e) => setOrderNumber(e.target.value)}
            className="mt-2 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div className="mb-4">
          <label htmlFor="orderDate" className="block text-sm font-medium text-gray-700">
            Start-Date
          </label>
          <input
            type="date"
            id="orderDate"
            value={orderDate}
            onChange={(e) => setOrderDate(e.target.value)}
            className="mt-2 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div className="mb-6">
          <label htmlFor="deliveryDate" className="block text-sm font-medium text-gray-700">
            End-Date
          </label>
          <input
            type="date"
            id="deliveryDate"
            value={deliveryDate}
            onChange={(e) => setDeliveryDate(e.target.value)}
            className="mt-2 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>


        <button
          onClick={downloadFile}
          disabled={loading || !orderDate || !deliveryDate}
          className={`w-full py-2 text-white font-semibold rounded-md shadow-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 ${loading || !orderDate || !deliveryDate
            ? "bg-gray-400 cursor-not-allowed"
            : "bg-indigo-600 hover:bg-indigo-700"
            }`}
        >
          {loading ? "Generating..." : "Generate Invoice"}
        </button>
      </div>
    </div>
  );
};

export default OrdersPage;
