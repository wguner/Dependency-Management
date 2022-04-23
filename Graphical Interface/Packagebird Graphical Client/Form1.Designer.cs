namespace Packagebird_Graphical_Client
{
    partial class Form1
    {
        /// <summary>
        ///  Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        ///  Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        ///  Required method for Designer support - do not modify
        ///  the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.button1 = new System.Windows.Forms.Button();
            this.button2 = new System.Windows.Forms.Button();
            this.label1 = new System.Windows.Forms.Label();
            this.registryPackagesList = new System.Windows.Forms.ListBox();
            this.registryProjectsList = new System.Windows.Forms.ListBox();
            this.label2 = new System.Windows.Forms.Label();
            this.button3 = new System.Windows.Forms.Button();
            this.newProjectNameBox = new System.Windows.Forms.TextBox();
            this.label3 = new System.Windows.Forms.Label();
            this.projectPackagesList = new System.Windows.Forms.ListBox();
            this.button4 = new System.Windows.Forms.Button();
            this.label4 = new System.Windows.Forms.Label();
            this.button5 = new System.Windows.Forms.Button();
            this.button6 = new System.Windows.Forms.Button();
            this.label5 = new System.Windows.Forms.Label();
            this.textBox2 = new System.Windows.Forms.TextBox();
            this.label6 = new System.Windows.Forms.Label();
            this.textBox3 = new System.Windows.Forms.TextBox();
            this.commandLineOutputTextbox = new System.Windows.Forms.RichTextBox();
            this.outputLabel = new System.Windows.Forms.Label();
            this.SuspendLayout();
            // 
            // button1
            // 
            this.button1.Location = new System.Drawing.Point(12, 301);
            this.button1.Name = "button1";
            this.button1.Size = new System.Drawing.Size(156, 23);
            this.button1.TabIndex = 0;
            this.button1.Text = "Add Package";
            this.button1.UseVisualStyleBackColor = true;
            // 
            // button2
            // 
            this.button2.Location = new System.Drawing.Point(328, 301);
            this.button2.Name = "button2";
            this.button2.Size = new System.Drawing.Size(156, 23);
            this.button2.TabIndex = 1;
            this.button2.Text = "Remove Package";
            this.button2.UseVisualStyleBackColor = true;
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(332, 18);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(99, 15);
            this.label1.TabIndex = 2;
            this.label1.Text = "Project Packages ";
            this.label1.Click += new System.EventHandler(this.label1_Click);
            // 
            // registryPackagesList
            // 
            this.registryPackagesList.FormattingEnabled = true;
            this.registryPackagesList.ItemHeight = 15;
            this.registryPackagesList.Location = new System.Drawing.Point(12, 36);
            this.registryPackagesList.Name = "registryPackagesList";
            this.registryPackagesList.Size = new System.Drawing.Size(156, 259);
            this.registryPackagesList.TabIndex = 3;
            // 
            // registryProjectsList
            // 
            this.registryProjectsList.FormattingEnabled = true;
            this.registryProjectsList.ItemHeight = 15;
            this.registryProjectsList.Location = new System.Drawing.Point(170, 36);
            this.registryProjectsList.Name = "registryProjectsList";
            this.registryProjectsList.Size = new System.Drawing.Size(156, 259);
            this.registryProjectsList.TabIndex = 4;
            // 
            // label2
            // 
            this.label2.AutoSize = true;
            this.label2.Location = new System.Drawing.Point(12, 18);
            this.label2.Name = "label2";
            this.label2.Size = new System.Drawing.Size(101, 15);
            this.label2.TabIndex = 5;
            this.label2.Text = "Registry Packages";
            this.label2.Click += new System.EventHandler(this.label2_Click);
            // 
            // button3
            // 
            this.button3.AccessibleName = "New Project";
            this.button3.Location = new System.Drawing.Point(170, 329);
            this.button3.Name = "button3";
            this.button3.Size = new System.Drawing.Size(156, 23);
            this.button3.TabIndex = 6;
            this.button3.Text = "New Project";
            this.button3.UseVisualStyleBackColor = true;
            this.button3.Click += new System.EventHandler(this.button3_Click);
            // 
            // newProjectNameBox
            // 
            this.newProjectNameBox.Location = new System.Drawing.Point(170, 387);
            this.newProjectNameBox.Name = "newProjectNameBox";
            this.newProjectNameBox.Size = new System.Drawing.Size(156, 23);
            this.newProjectNameBox.TabIndex = 7;
            // 
            // label3
            // 
            this.label3.AutoSize = true;
            this.label3.Location = new System.Drawing.Point(82, 390);
            this.label3.Name = "label3";
            this.label3.Size = new System.Drawing.Size(82, 15);
            this.label3.TabIndex = 8;
            this.label3.Text = "Project Name:";
            // 
            // projectPackagesList
            // 
            this.projectPackagesList.FormattingEnabled = true;
            this.projectPackagesList.ItemHeight = 15;
            this.projectPackagesList.Location = new System.Drawing.Point(332, 36);
            this.projectPackagesList.Name = "projectPackagesList";
            this.projectPackagesList.Size = new System.Drawing.Size(152, 259);
            this.projectPackagesList.TabIndex = 9;
            // 
            // button4
            // 
            this.button4.Location = new System.Drawing.Point(170, 301);
            this.button4.Name = "button4";
            this.button4.Size = new System.Drawing.Size(156, 23);
            this.button4.TabIndex = 10;
            this.button4.Text = "Install Project";
            this.button4.UseVisualStyleBackColor = true;
            // 
            // label4
            // 
            this.label4.AutoSize = true;
            this.label4.Location = new System.Drawing.Point(170, 18);
            this.label4.Name = "label4";
            this.label4.Size = new System.Drawing.Size(94, 15);
            this.label4.TabIndex = 11;
            this.label4.Text = "Registry Projects";
            // 
            // button5
            // 
            this.button5.Location = new System.Drawing.Point(170, 358);
            this.button5.Name = "button5";
            this.button5.Size = new System.Drawing.Size(156, 23);
            this.button5.TabIndex = 12;
            this.button5.Text = "Sync Project";
            this.button5.UseVisualStyleBackColor = true;
            // 
            // button6
            // 
            this.button6.Location = new System.Drawing.Point(328, 329);
            this.button6.Name = "button6";
            this.button6.Size = new System.Drawing.Size(156, 23);
            this.button6.TabIndex = 13;
            this.button6.Text = "New Package";
            this.button6.UseVisualStyleBackColor = true;
            this.button6.Click += new System.EventHandler(this.button6_Click);
            // 
            // label5
            // 
            this.label5.AutoSize = true;
            this.label5.Location = new System.Drawing.Point(486, 39);
            this.label5.Name = "label5";
            this.label5.Size = new System.Drawing.Size(74, 15);
            this.label5.TabIndex = 14;
            this.label5.Text = "Current Path";
            this.label5.Click += new System.EventHandler(this.label5_Click);
            // 
            // textBox2
            // 
            this.textBox2.Location = new System.Drawing.Point(566, 36);
            this.textBox2.Name = "textBox2";
            this.textBox2.Size = new System.Drawing.Size(214, 23);
            this.textBox2.TabIndex = 15;
            // 
            // label6
            // 
            this.label6.AutoSize = true;
            this.label6.Location = new System.Drawing.Point(486, 67);
            this.label6.Name = "label6";
            this.label6.Size = new System.Drawing.Size(73, 15);
            this.label6.TabIndex = 16;
            this.label6.Text = "Current User";
            // 
            // textBox3
            // 
            this.textBox3.Location = new System.Drawing.Point(566, 64);
            this.textBox3.Name = "textBox3";
            this.textBox3.Size = new System.Drawing.Size(214, 23);
            this.textBox3.TabIndex = 17;
            // 
            // commandLineOutputTextbox
            // 
            this.commandLineOutputTextbox.Location = new System.Drawing.Point(490, 199);
            this.commandLineOutputTextbox.Name = "commandLineOutputTextbox";
            this.commandLineOutputTextbox.Size = new System.Drawing.Size(298, 96);
            this.commandLineOutputTextbox.TabIndex = 18;
            this.commandLineOutputTextbox.Text = "";
            // 
            // outputLabel
            // 
            this.outputLabel.AutoSize = true;
            this.outputLabel.Location = new System.Drawing.Point(490, 181);
            this.outputLabel.Name = "outputLabel";
            this.outputLabel.Size = new System.Drawing.Size(133, 15);
            this.outputLabel.TabIndex = 19;
            this.outputLabel.Text = "Command Line Output:";
            this.outputLabel.Click += new System.EventHandler(this.label7_Click);
            // 
            // Form1
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(7F, 15F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(800, 450);
            this.Controls.Add(this.outputLabel);
            this.Controls.Add(this.commandLineOutputTextbox);
            this.Controls.Add(this.textBox3);
            this.Controls.Add(this.label6);
            this.Controls.Add(this.textBox2);
            this.Controls.Add(this.label5);
            this.Controls.Add(this.button6);
            this.Controls.Add(this.button5);
            this.Controls.Add(this.label4);
            this.Controls.Add(this.button4);
            this.Controls.Add(this.projectPackagesList);
            this.Controls.Add(this.label3);
            this.Controls.Add(this.newProjectNameBox);
            this.Controls.Add(this.button3);
            this.Controls.Add(this.label2);
            this.Controls.Add(this.registryProjectsList);
            this.Controls.Add(this.registryPackagesList);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.button2);
            this.Controls.Add(this.button1);
            this.Name = "Form1";
            this.Text = "Form1";
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private Button button1;
        private Button button2;
        private Label label1;
        private ListBox registryPackagesList;
        private ListBox registryProjectsList;
        private Label label2;
        private Button button3;
        private TextBox newProjectNameBox;
        private Label label3;
        private ListBox projectPackagesList;
        private Button button4;
        private Label label4;
        private Button button5;
        private Button button6;
        private Label label5;
        private TextBox textBox2;
        private Label label6;
        private TextBox textBox3;
        private RichTextBox commandLineOutputTextbox;
        private Label outputLabel;
    }
}