import React from "react";

const LogisticsPerformaceTable = () => {

    const data = [
        {
          location: 'Wellsville',
          name: 'Cami',
          orderCount: 29,
          revenue: 56164.79,
          truckHire: 45260.52,
          net: 10904.27,
          marginPercentage: 19.41,
          miles: 21295,
          revenueMile: 0.51,
          orderOTP: 71.43,
          stopOTP: 86.96,
          revOrderDay: 1936.72,
          netOrderDay: 376.01,
          ordersDay: 5,
        },
        {
          location: 'Wellsville',
          name: 'Lenora',
          orderCount: 41,
          revenue: 74011.2,
          truckHire: 58120.35,
          net: 15890.85,
          marginPercentage: 21.47,
          miles: 25566,
          revenueMile: 0.62,
          orderOTP: 88.89,
          stopOTP: 93.24,
          revOrderDay: 1805.15,
          netOrderDay: 387.58,
          ordersDay: 7,
        },
        // ... (rest of the data)
      ];
    
      return (
        <div>
        <div className="flex flex-row gap-8 space-betwen">
          <div className="overflow-x-auto mb-8">
            <h2 className="text-xl font-bold mb-4">Place</h2>
            <table className="table-auto">
              <thead>
                <tr className="bg-gray-800 text-white">
                  <th className="px-4 py-2"></th>
                  <th className="px-4 py-2">ORDER COUNT</th>
                  <th className="px-4 py-2">REVENUE</th>
                  <th className="px-4 py-2">TRUCK HIRE</th>
                  <th className="px-4 py-2">NET</th>
                  <th className="px-4 py-2">MARGIN%</th>
                  <th className="px-4 py-2">MILES</th>
                  <th className="px-4 py-2">REV/MILE</th>
                  <th className="px-4 py-2">ORDER</th>
                  <th className="px-4 py-2">STOP</th>
                </tr>
              </thead>
              <tbody>
                {data.map((item, index) => (
                  <tr key={index} className={index % 2 === 0 ? 'bg-gray-100' : 'bg-white'}>
                    <td className="border px-4 py-2">{item.name}</td>
                    <td className="border px-4 py-2 text-right">{item.orderCount}</td>
                    <td className="border px-4 py-2 text-right">{item.revenue.toFixed(2)}</td>
                    <td className="border px-4 py-2 text-right">{item.truckHire.toFixed(2)}</td>
                    <td className="border px-4 py-2 text-right">{item.net.toFixed(2)}</td>
                    <td className="border px-4 py-2 text-right">{item.marginPercentage.toFixed(2)}%</td>
                    <td className="border px-4 py-2 text-right">{item.miles}</td>
                    <td className="border px-4 py-2 text-right">{item.revenueMile.toFixed(2)}</td>
                    <td className="border px-4 py-2 text-right">{item.orderOTP.toFixed(2)}%</td>
                    <td className="border px-4 py-2 text-right">{item.stopOTP.toFixed(2)}%</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

    
          <div className="overflow-x-auto">
            <h2 className="text-xl font-bold mb-4">Metrics</h2>
            <table className="table-auto">
              <thead>
                <tr className="bg-gray-800 text-white">
                  <th className="px-4 py-2">LOCATION</th>
                  <th className="px-4 py-2">REV/ORDER/DAY</th>
                  <th className="px-4 py-2">NET/ORDER/DAY</th>
                  <th className="px-4 py-2">ORDERS/DAY</th>
                </tr>
              </thead>
              <tbody>
                {data.map((item, index) => (
                  <tr key={index} className={index % 2 === 0 ? 'bg-gray-100' : 'bg-white'}>
                    <td className="border px-4 py-2">{item.location}</td>
                    <td className="border px-4 py-2 text-right">{item.revOrderDay.toFixed(2)}</td>
                    <td className="border px-4 py-2 text-right">{item.netOrderDay.toFixed(2)}</td>
                    <td className="border px-4 py-2 text-right">{item.ordersDay}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
        </div>
      );
    };
    
export default LogisticsPerformaceTable;