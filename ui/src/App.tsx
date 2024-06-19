import Dashboard from "./layout/Dashboard";
import Header from "./layout/Header";

function App() {
  return (
    <div className="h-dvh flex flex-col pb-6 gap-6">
      <Header />
      <Dashboard />
    </div>
  );
}

export default App;
