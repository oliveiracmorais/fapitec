"use client";

import { formatarStatus, corStatus } from "../lib/api-propostas";

type Props = {
  status: string;
};

export default function PropostaStatusBadge({ status }: Props) {
  return (
    <span
      className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${corStatus(status)}`}
    >
      {formatarStatus(status)}
    </span>
  );
}
