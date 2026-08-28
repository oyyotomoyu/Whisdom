import "./Spinner.css";

interface SpinnerProps {
  fullscreen?: boolean;
  size?: number;
}

export function Spinner({ fullscreen, size = 20 }: SpinnerProps) {
  const spinner = (
    <span
      className="spinner"
      style={{ width: size, height: size }}
      role="status"
      aria-label="Loading"
    />
  );

  if (!fullscreen) return spinner;

  return <div className="spinner-fullscreen">{spinner}</div>;
}
