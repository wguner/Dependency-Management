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

            var loadPackages = GenerateCommand($"get packages");
            line = "";

            loadPackages.Start();
            while (true)
            {
                char i = (char)loadPackages.StandardOutput!.Read();
                if (i.Equals('\n'))
                {
                    this.registryPackagesList.Items.Add(line);
                    line = "";
                }
                else
                {
                    line += i;
                }
                if (loadPackages.StandardOutput?.EndOfStream == true)
                    break;
            }
            loadPackages.WaitForExit();

            var outMsg = loadPackages.StandardOutput.ReadToEnd();
            var errMsg = loadPackages.StandardError.ReadToEnd();

            this.commandLineOutputTextbox.Text = $"StdOut: {outMsg}\nStdErr: {errMsg}";

            // Add loading project packages tomorrow
        }

        private void label1_Click(object sender, EventArgs e)
        {

        }

        private void label2_Click(object sender, EventArgs e)
        {

        }

        private void button6_Click(object sender, EventArgs e)
        {
            var project = registryProjectsList.SelectedItem.ToString();
            var process = GenerateCommand($"create package {project} 0");
            process.StartInfo.WorkingDirectory = $".\\{project}";
            process.Start();
            registryPackagesList.Items.Clear();
            process = GenerateCommand("get packages");
            process.Start();
            string line = "";
            while (true)
            {
                char i = (char)process.StandardOutput!.Read();
                if (i.Equals('\n'))
                {
                    this.registryPackagesList.Items.Add(line);
                    line = "";
                }
                else
                {
                    line += i;
                }
                if (process.StandardOutput?.EndOfStream == true)
                    break;
            }
            process.WaitForExit();
            var outMsg = process.StandardOutput.ReadToEnd();
            var errMsg = process.StandardError.ReadToEnd();

            this.commandLineOutputTextbox.Text = $"StdOut: {outMsg}\nStdErr: {errMsg}";
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
            this.registryProjectsList.Items.Clear();

            var getProjects = GenerateCommand($"get projects");
            getProjects.Start();
            string line = "";

            while (true)
            {
                char i = (char)getProjects.StandardOutput.Read();
                if (i.Equals('\n'))
                {
                    this.registryProjectsList.Items.Add(line);
                    line = "";
                }
                else
                {
                    line += i;
                }
                if (getProjects.StandardOutput.EndOfStream == true)
                    break;
            }
            getProjects.WaitForExit();
        }

        private void label7_Click(object sender, EventArgs e)
        {

        }

        private void projectPackagesList_SelectedIndexChanged(object sender, EventArgs e)
        {

        }

        private void button1_Click(object sender, EventArgs e)
        {
            var selectedPackage = this.registryPackagesList.SelectedItem.ToString();
            string selectedPackageName = selectedPackage!.Split("-v")[0], selectedPackageVersion = selectedPackage!.Split("-v")[1];
            var addPackage = GenerateCommand($"add package {selectedPackageName} {selectedPackageVersion}");
            addPackage.StartInfo.WorkingDirectory = $".\\{registryProjectsList.SelectedItem.ToString()}";
            addPackage.Start();
            string outMsg = addPackage.StandardOutput.ReadToEnd(), errMsg = addPackage.StandardError.ReadToEnd();
            addPackage.WaitForExitAsync();
            var getProjectPackages = GenerateCommand($"get added");
            getProjectPackages.StartInfo.WorkingDirectory = $".\\{registryProjectsList.SelectedItem.ToString()}";
            getProjectPackages.Start();
            string line = "";
            char c;
            while (true)
            {
                c = (char)getProjectPackages.StandardOutput.Read();
                if (c == '\n')
                {
                    this.projectPackagesList.Items.Add(line);
                } else
                {
                    line += c;
                }
                if (getProjectPackages.StandardOutput.EndOfStream)
                    break;
            }
            this.commandLineOutputTextbox.Text = $"Package: {selectedPackageName} Version: {selectedPackageVersion}\nStdOut: {outMsg}\nStdErr: {errMsg}";
            getProjectPackages.WaitForExit();
        }

        private void button5_Click(object sender, EventArgs e)
        {
            var syncProject = GenerateCommand("sync");
            syncProject.StartInfo.WorkingDirectory = $".\\{registryProjectsList.SelectedItem.ToString()}";
            syncProject.Start();
            syncProject.WaitForExitAsync();
        }

        private void registryProjectsList_SelectedIndexChanged(object sender, EventArgs e)
        {
            var getProjectPackages = GenerateCommand($"get added");
            if (this.registryProjectsList.Equals(null) || this.registryProjectsList.SelectedIndex.Equals(-1))
                return;
            this.projectPackagesList.Items.Clear();
            getProjectPackages.StartInfo.WorkingDirectory = $".\\{registryProjectsList.SelectedItem.ToString()}";
            getProjectPackages.Start();
            string line = "";
            char c;
            while (true)
            {
                c = (char)getProjectPackages.StandardOutput.Read();
                if (c == '\n')
                {
                    this.projectPackagesList.Items.Add(line);
                }
                else
                {
                    line += c;
                }
                if (getProjectPackages.StandardOutput.EndOfStream)
                    break;
            }
            getProjectPackages.WaitForExit();
        }
    }
}