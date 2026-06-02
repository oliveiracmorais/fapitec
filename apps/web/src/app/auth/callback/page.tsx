import { Suspense } from "react";
import CallbackContent from "./content";

export default function CallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-brand-900 via-brand-800 to-brand-700">
          <div className="rounded-xl bg-white p-8 shadow-2xl">
            <p className="text-gray-700">Autenticando...</p>
          </div>
        </div>
      }
    >
      <CallbackContent />
    </Suspense>
  );
}
