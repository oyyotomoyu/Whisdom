import { Link } from "react-router-dom";
import { EmptyState } from "../../components/EmptyState";

export function NotFound() {
  return (
    <div style={{ height: "100%" }}>
      <EmptyState
        title="Page not found"
        subtitle="The page you are looking for doesn't exist."
        action={<Link to="/chat">Back to chat</Link>}
      />
    </div>
  );
}
