import "./summaryTable.scss";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";

const rows = [
  {
    id: 1000003,
    customer: "Joe",
    time: "14:40",
    date: "12 May 2012",
    amount: 414.18,
    paymentMethod: "card",
    driver: "Lawrence",
    status: "Pending",
  },
  {
    id: 1000002,
    customer: "Bobby Singer",
    time: "14:25",
    date: "12 May 2012",
    amount: 12.48,
    paymentMethod: "cash",
    driver: "John Winchester",
    status: "OnGoing",
  },
  {
    id: 1000001,
    customer: "Crowley",
    time: "00:25",
    date: "10 May 2012",
    amount: 666.66,
    paymentMethod: "online",
    driver: "Michael",
    status: "Completed",
  },
  {
    id: 1000000,
    customer: "Daniel",
    time: "05:41",
    date: "8 May 2012",
    amount: 592.14,
    paymentMethod: "cash",
    driver: "Richard",
    status: "Canceled",
  },
];

const SummaryTable = () => {
  return (
    <div className="table">
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell className="tableCell">Customer</TableCell>
              <TableCell className="tableCell">Time & Date</TableCell>
              <TableCell className="tableCell">Status</TableCell>
              <TableCell className="tableCell">Driver</TableCell>
              <TableCell className="tableCell">Amount (RM)</TableCell>
              <TableCell className="tableCell">Payment Method</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row) => (
              <TableRow
                key={row.id}
                sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
              >
                <TableCell component="th" scope="row">
                  {row.id}
                </TableCell>
                <TableCell className="tableCell">{row.customer}</TableCell>
                <TableCell className="tableCell">
                  {row.time} {row.date}
                </TableCell>
                <TableCell className="tableCell">
                  <span className={`status ${row.status}`}>{row.status}</span>
                </TableCell>
                <TableCell className="tableCell">{row.driver}</TableCell>
                <TableCell className="tableCell">{row.amount}</TableCell>
                <TableCell className="tableCell">{row.paymentMethod}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default SummaryTable;
