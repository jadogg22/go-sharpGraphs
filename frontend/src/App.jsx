import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, NavLink } from 'react-router-dom';
import DailyOps from './DailyOps';
import LogisticsWeekly from './LogisticsWeekly';

const App = () => {
  return (
    <Router>
      <div className="flex h-screen">
        <div className="w-64 bg-gray-800 text-white p-4 flex flex-col">
          <nav className="flex-grow flex flex-col justify-center">
            <ul>
              <li>
                <NavLink
                  end
                  to="/"
                  activeClassName="p-4 bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Home
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/dailyOps"
                  activeClassName="bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Daily
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/about"
                  activeClassName="bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  about
                </NavLink>
              </li>
              <li>
                <NavLink 
                to="/logisticsWeekly"
                activeClassName="bg-gray-700"
                className="block py-2 px-4 hover:bg-gray-700"
                >
                  Logistics Weekly
                </NavLink>
              </li>
            </ul>
          </nav>
        </div>
        <div className="flex-1 p-6 flex flex-col justify-center overflow-auto">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/dailyOps" element={<DailyOps />} />
            <Route path="/about" element={<About />} />
            <Route path="/logisticsWeekly" element={<LogisticsWeekly />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
};

const Home = () => (
  <div>
    <h1 className="text-2xl font-bold">Home Page</h1>
    <p>Welcome to the home page!</p>
  </div>
);

const About = () => (
  <div>
    <h1 className="text-2xl font-bold">About Page</h1>
    <p>Learn more about us on the about page.</p>
  </div>
);


export default App;