"use client";

import { useRouter } from "next/navigation";
import FormAvaliador from "../../../../components/avaliadores/FormAvaliador";
import { cadastrarAvaliador } from "../../../../lib/api-avaliadores";

export default function NovoAvaliadorPage() {
  const router = useRouter();

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  async function handleSubmit(data: any) {
    await cadastrarAvaliador(data);
    router.push("/avaliadores");
  }

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <FormAvaliador
        mode="criar"
        onSubmit={handleSubmit}
        onCancel={() => router.push("/avaliadores")}
      />
    </main>
  );
}
