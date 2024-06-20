import CodeIcon from '../../assets/undraw_code.svg?react';

export default function Help() {
  const url = new URL(window.location.href);

  return (
    <div className="flex w-full flex-col items-center gap-4 pt-8">
      <CodeIcon className="h-40 w-auto" />
      <div className="mt-8">Integrate the following script into your web application to record user sessions.</div>
      <div className="rounded-md bg-gray-200 px-4 py-2 font-mono text-sm">{`<script src="${url.origin}/cassette.min.cjs"></script`}</div>
    </div>
  );
}
