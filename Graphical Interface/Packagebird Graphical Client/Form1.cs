using System;
using System.Diagnostics;
namespace Packagebird_Graphical_Client
{
    public partial class Form1 : Form
    {
        const string executablePath = "C:\\Users\\Elisha Aguilera\\GolandProjects\\Dependency-Management-Beta-Update\\builds\\client\\packagebird.exe";

        public Process GenerateCommand(string arguments)
        {
            Process process = new Process
            {
                StartInfo = new ProcessStartInfo
                {
                    FileName = executablePath,
                    Arguments = arguments,
                    UseShellExecute = false,
                    RedirectStandardError = true,
                    RedirectStandardOutput = true,
                    CreateNoWindow = true,
                }
            };
            return process;
        }

        public Form1()
        {
            InitializeComponent();
            Process loadProjects = new Process
            {
                StartInfo = new ProcessStartInfo
                {
                    FileName = executablePath,
                    Arguments = "get projects",
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    CreateNoWindow = true,
                }
            };
            loadProjects.Start();
            string line = "";

            while (true) {
                char i = (char)loadProjects.StandardOutput.Read();
                if (i.Equals('\n'))
                {
                    this.registryProjectsList.Items.Add(line);
                    line = "";
                } else
                {
                    line += i;
                }
                if (loadProjects.StandardOutput.EndOfStream == true)
                    break;
            }
            loadProjects.WaitForExit();
        }

        private void label1_Click(object sender, EventArgs e)
        {

        }

        private void label2_Click(object sender, EventArgs e)
        {

        }

        private void button6_Click(object sender, EventArgs e)
        {

        }

        private void label5_Click(object sender, EventArgs e)
        {

        }

        private void listView1_SelectedIndexChanged(object sender, EventArgs e)
        {

        }

        private void button3_Click(object sender, EventArgs e)
        {
            string projectName = this.newProjectNameBox.Text;
            
            if (string.IsNullOrEmpty(projectName))
                return;

            Process newProject = GenerateCommand($"create project {projectName}");
            newProject.Start();
            newProject.WaitForExit();

            var outMsg = newProject.StandardOutput.ReadToEnd();
            var errMsg = newProject.StandardError.ReadToEnd();

            this.commandLineOutputTextbox.Text = $"StdOut: {outMsg}\nStdErr: {errMsg}";
        }

        private void label7_Click(object sender, EventArgs e)
        {

        }
    }
}