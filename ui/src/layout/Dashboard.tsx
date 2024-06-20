import { Outlet } from 'react-router-dom';
import Header from './Header';

export default function Dashboard() {
  return (
    <div className="flex h-dvh flex-col gap-6 pb-6">
      <Header />
      <div className="mx-auto flex h-screen w-full max-w-7xl gap-4">
        <Outlet />
      </div>
    </div>
  );
}
