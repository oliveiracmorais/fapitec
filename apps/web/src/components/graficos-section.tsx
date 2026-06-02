"use client";

import { useEffect, useState } from "react";
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Legend,
} from "recharts";

type Faixa = {
  nome: string;
  valor: number;
  legenda: string;
};

type DadoLinha = {
  mes: string;
  valor: number;
};

type GraficoDonut = {
  id: string;
  tipo: "donut";
  titulo: string;
  faixas: Faixa[];
};

type GraficoLinha = {
  id: string;
  tipo: "linha";
  titulo: string;
  dados: DadoLinha[];
};

type Grafico = GraficoDonut | GraficoLinha;

const CORES_DONUT = ["#7C3AED", "#10B981", "#F59E0B"];

function CustomTooltip({ active, payload, label }: any) {
  if (active && payload && payload.length > 0) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white px-3 py-2 shadow-lg">
        <p className="text-sm font-medium text-gray-900">
          {label || payload[0].name}
        </p>
        <p className="text-sm text-gray-600">{`${payload[0].value}%`}</p>
      </div>
    );
  }
  return null;
}

export default function GraficosSection() {
  const [graficos, setGraficos] = useState<Grafico[]>([]);
  const [carregando, setCarregando] = useState(true);

  useEffect(() => {
    fetch("/api/v1/dashboard/graficos")
      .then((res) => res.json())
      .then((data) => setGraficos(data.graficos || []))
      .catch(() => setGraficos([]))
      .finally(() => setCarregando(false));
  }, []);

  if (carregando) {
    return (
      <section aria-label="Gráficos" className="mb-8">
        <h2 className="mb-4 text-lg font-bold text-gray-900">Gráficos</h2>
        <p className="text-sm text-gray-500">Carregando gráficos...</p>
      </section>
    );
  }

  if (graficos.length === 0) return null;

  return (
    <section aria-label="Gráficos do sistema" className="mb-8">
      <h2 className="mb-4 text-lg font-bold text-gray-900">Gráficos</h2>
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {graficos.map((grafico) => {
          if (grafico.tipo === "donut") {
            return (
                <div
                  key={grafico.id}
                  className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm"
                  title={grafico.titulo}
                  role="region"
                  aria-label={`Gráfico: ${grafico.titulo}`}
                >
                  <h3 className="mb-4 text-sm font-semibold text-gray-700">
                    {grafico.titulo}
                  </h3>
                  <ResponsiveContainer width="100%" height={260}>
                    <PieChart>
                      <Pie
                        data={grafico.faixas}
                        dataKey="valor"
                        nameKey="nome"
                        cx="50%"
                        cy="50%"
                        innerRadius={60}
                        outerRadius={100}
                        paddingAngle={3}
                      >
                        {grafico.faixas.map((_, i) => (
                          <Cell
                            key={i}
                            fill={CORES_DONUT[i % CORES_DONUT.length]}
                          />
                        ))}
                      </Pie>
                      <Tooltip content={<CustomTooltip />} />
                      <Legend
                        formatter={(value: string) => (
                          <span className="text-sm text-gray-600">{value}</span>
                        )}
                      />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
            );
          }

          if (grafico.tipo === "linha") {
            return (
              <div
                  key={grafico.id}
                  className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm"
                  title={grafico.titulo}
                  role="region"
                  aria-label={`Gráfico: ${grafico.titulo}`}
                >
                  <h3 className="mb-4 text-sm font-semibold text-gray-700">
                    {grafico.titulo}
                  </h3>
                  <ResponsiveContainer width="100%" height={260}>
                    <LineChart data={grafico.dados}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
                    <XAxis
                      dataKey="mes"
                      tick={{ fontSize: 12, fill: "#6B7280" }}
                      axisLine={{ stroke: "#E5E7EB" }}
                    />
                    <YAxis
                      tick={{ fontSize: 12, fill: "#6B7280" }}
                      axisLine={{ stroke: "#E5E7EB" }}
                    />
                    <Tooltip content={<CustomTooltip />} />
                    <Line
                      type="monotone"
                      dataKey="valor"
                      stroke="#7C3AED"
                      strokeWidth={2}
                      dot={{ fill: "#7C3AED", r: 4 }}
                      activeDot={{ r: 6 }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            );
          }

          return null;
        })}
      </div>
    </section>
  );
}
