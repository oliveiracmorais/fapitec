import Header from "../../components/header";
import Breadcrumb from "../../components/breadcrumb";

export default function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <Breadcrumb />
      {children}
    </div>
  );
}
