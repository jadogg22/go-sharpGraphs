import React from 'react';

const PerformanceTable1 = () => {
  const data = [
    {
      driverManager: 'Felicia McMichael',
      numTrucks: 15,
      milesPerTruck: 1036,
      deadheadPercentage: '2.40%',
      orderPercentage: '100.00%',
      stopPercentage: '100.00%',
    },
    {
      driverManager: 'Katrina Sepulveda',
      numTrucks: 17,
      milesPerTruck: 818,
      deadheadPercentage: '2.40%',
      orderPercentage: '100.00%',
      stopPercentage: '100.00%',
    },
    {
      driverManager: 'Ken Moss JR.',
      numTrucks: 17,
      milesPerTruck: 752,
      deadheadPercentage: '11.00%',
      orderPercentage: '100.00%',
      stopPercentage: '100.00%',
    },
    {
      driverManager: 'Lindsay Workman',
      numTrucks: 15,
      milesPerTruck: 574,
      deadheadPercentage: '5.40%',
      orderPercentage: '100.00%',
      stopPercentage: '100.00%',
    },
    {
      driverManager: 'Tracy Rigby',
      numTrucks: 14,
      milesPerTruck: 787,
      deadheadPercentage: '4.70%',
      orderPercentage: '100.00%',
      stopPercentage: '100.00%',
    },
    {
      driverManager: 'OTR Totals',
      numTrucks: 78,
      milesPerTruck: 793,
      deadheadPercentage: '5.18%',
      orderPercentage: '100.00%',
      stopPercentage: '100.00%',
    },
    {
      driverManager: 'Rochelle Genera',
      numTrucks: 21,
      milesPerTruck: 150,
      deadheadPercentage: '29.90%',
      orderPercentage: '88.99%',
      stopPercentage: '95.00%',
    },
    {
      driverManager: 'Stephanie Bingham',
      numTrucks: 3,
      milesPerTruck: 203,
      deadheadPercentage: '28.20%',
      orderPercentage: '100.00%',
      stopPercentage: '100.00%',
    },
  ];

  return (
    <div className="overflow-x-auto">
      <table className="table-auto">
        <thead>
          <tr className="bg-gray-800 text-white">
            <th className="px-4 py-2">Driver Manager</th>
            <th className="px-4 py-2"># Trucks</th>
            <th className="px-4 py-2">Miles Per Truck</th>
            <th className="px-4 py-2">Deadhead %</th>
            <th className="px-4 py-2">ORDER</th>
            <th className="px-4 py-2">STOP</th>
          </tr>
        </thead>
        <tbody>
          {data.map((item, index) => (
            <tr key={index} className={index % 2 === 0 ? 'bg-gray-100' : 'bg-white'}>
              <td className="border px-4 py-2">{item.driverManager}</td>
              <td className="border px-4 py-2 text-right">{item.numTrucks}</td>
              <td className="border px-4 py-2 text-right">{item.milesPerTruck}</td>
              <td className="border px-4 py-2 text-right">{item.deadheadPercentage}</td>
              <td className="border px-4 py-2 text-right">{item.orderPercentage}</td>
              <td className="border px-4 py-2 text-right">{item.stopPercentage}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default PerformanceTable1;