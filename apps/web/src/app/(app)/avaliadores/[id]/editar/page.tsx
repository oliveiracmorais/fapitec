"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import FormAvaliador from "../../../../../components/avaliadores/FormAvaliador";
import { visualizarAvaliador, editarAvaliador } from "../../../../../lib/api-avaliadores";
import type { AvaliadorSaida } from "../../../../../lib/api-avaliadores";

export default function EditarAvaliadorPage() {
  const router = useRouter();
  const params = useParams();
  const [avaliador, setAvaliador] = useState<AvaliadorSaida | null>(null);
  const [carregando, setCarregando] = useState(true);

  useEffect(() => {
    if (!params?.id) return;
    visualizarAvaliador(Number(params.id))
      .then(setAvaliador)
      .catch(() => router.push("/avaliadores"))
      .finally(() => setCarregando(false));
  }, [params, router]);

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  async function handleSubmit(data: any) {
    if (!params?.id) return;
    await editarAvaliador(Number(params.id), data);
    router.push(`/avaliadores/${params.id}`);
  }

  if (carregando) {
    return (
      <div className="flex items-center justify-center py-20 text-gray-600">Carregando...</div>
    );
  }

  if (!avaliador) return null;

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <FormAvaliador
        mode="editar"
        initialData={avaliador}
        onSubmit={handleSubmit}
        onCancel={() => router.push(`/avaliadores/${params!.id}`)}
      />
    </main>
  );
}
