import { useState } from "react";

export default function VacationExport() {
  const [loading, setLoading] = useState(false);

  const downloadCSV = async () => {
    setLoading(true);

    try {
      const res = await fetch("api2/vacation/all");
      if (!res.ok) throw new Error("Failed to fetch data");
      const json = await res.json();

      // Flatten JSON into CSV rows
      const headers = [
        "Company",
        "EmployeeID",
        "EmployeeName",
        "VacationHoursDue",
        "VacationHoursRate",
        "AmountDue",
      ];

      const rows = [];
      for (const [company, employees] of Object.entries(json.Data)) {
        employees.forEach((emp) => {
          rows.push([
            company,
            emp.EmployeeID,
            emp.EmployeeName,
            emp.VacationHoursDue,
            emp.VacationHoursRate,
            emp.AmountDue,
          ]);
        });
      }

      // Build CSV string
      const csvContent = [headers, ...rows]
        .map((row) => row.map((v) => `"${v}"`).join(","))
        .join("\n");

      // Trigger file download
      const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.setAttribute("href", url);
      link.setAttribute("download", "vacation.csv");
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

    } catch (err) {
      console.error("Error exporting CSV:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <h1 className="text-3xl text-center text-gray-800 mb-16">Vacation Export</h1>
      <div className="border-4 border-gray-600 bg-white p-4 shadow-lg rounded-lg flex flex-col items-center">
        <p className="text-gray-600 mb-4 text-center">
          Click the button below to download a CSV file containing all employee vacation data.
        </p>
        <button
          onClick={downloadCSV}
          disabled={loading}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg shadow hover:bg-blue-700 disabled:bg-gray-400 flex items-center"
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clipRule="evenodd" />
          </svg>
          {loading ? "Downloading..." : "Download CSV"}
        </button>
      </div>
    </div>
  );
}
